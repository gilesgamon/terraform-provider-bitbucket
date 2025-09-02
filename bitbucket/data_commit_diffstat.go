package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataCommitDiffstat() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCommitDiffstatRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commit_sha": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Commit SHA to retrieve diff statistics for",
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to filter diff statistics by",
			},
			"diffstat": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"old_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"new_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lines_added": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lines_removed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lines_changed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"binary": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"renamed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"new_file": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"summary": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataCommitDiffstatRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	commitSha := d.Get("commit_sha").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitDiffstatRead", dumpResourceData(d, dataCommitDiffstat().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/diffstat/%s", workspace, repoSlug, commitSha)

	// Build query parameters
	params := make(map[string]string)
	if path, ok := d.GetOk("path"); ok {
		params["path"] = path.(string)
	}

	// Add query parameters to URL
	if len(params) > 0 {
		url += "?"
		first := true
		for key, value := range params {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", key, value)
			first = false
		}
	}

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from commit diffstat call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commitSha, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading commit diffstat with params (%s): ", dumpResourceData(d, dataCommitDiffstat().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	diffstatBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] commit diffstat response: %v", diffstatBody)

	var diffstatResponse CommitDiffstatResponse
	decodeerr := json.Unmarshal(diffstatBody, &diffstatResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/diffstat", workspace, repoSlug, commitSha))
	flattenCommitDiffstat(&diffstatResponse, d)
	return nil
}

// CommitDiffstatResponse represents the response from the commit diffstat API
type CommitDiffstatResponse struct {
	Values []DiffstatFile `json:"values"`
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Next   string         `json:"next"`
}

// DiffstatFile represents a file in the commit diffstat
type DiffstatFile struct {
	Type        string                 `json:"type"`
	OldPath     string                 `json:"old_path"`
	NewPath     string                 `json:"new_path"`
	Status      string                 `json:"status"`
	LinesAdded  int                    `json:"lines_added"`
	LinesRemoved int                    `json:"lines_removed"`
	LinesChanged int                    `json:"lines_changed"`
	Size        int                    `json:"size"`
	Binary      bool                   `json:"binary"`
	Renamed     bool                   `json:"renamed"`
	Deleted     bool                   `json:"deleted"`
	NewFile     bool                   `json:"new_file"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the commit diffstat information
func flattenCommitDiffstat(c *CommitDiffstatResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	diffstat := make([]interface{}, len(c.Values))
	totalAdded := 0
	totalRemoved := 0
	totalChanged := 0

	for i, file := range c.Values {
		diffstat[i] = map[string]interface{}{
			"type":         file.Type,
			"old_path":     file.OldPath,
			"new_path":     file.NewPath,
			"status":       file.Status,
			"lines_added":  file.LinesAdded,
			"lines_removed": file.LinesRemoved,
			"lines_changed": file.LinesChanged,
			"size":         file.Size,
			"binary":       file.Binary,
			"renamed":      file.Renamed,
			"deleted":      file.Deleted,
			"new_file":     file.NewFile,
		}

		totalAdded += file.LinesAdded
		totalRemoved += file.LinesRemoved
		totalChanged += file.LinesChanged
	}

	// Set summary statistics
	summary := map[string]interface{}{
		"total_files":    len(c.Values),
		"total_added":    totalAdded,
		"total_removed":  totalRemoved,
		"total_changed":  totalChanged,
		"total_modified": totalAdded + totalRemoved,
	}

	d.Set("diffstat", diffstat)
	d.Set("summary", summary)
}
