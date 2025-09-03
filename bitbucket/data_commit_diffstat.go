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
			"commit": {
				Type:     schema.TypeString,
				Required: true,
			},
			"diffstat": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"new_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"old_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"new_file": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"renamed_file": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"deleted_file": {
							Type:     schema.TypeBool,
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataCommitDiffstatRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	commit := d.Get("commit").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitDiffstatRead", dumpResourceData(d, dataCommitDiffstat().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commits/%s/diffstat", workspace, repoSlug, commit)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from commit diffstat call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commit, workspace, repoSlug)
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

	d.SetId(fmt.Sprintf("%s/%s/commits/%s/diffstat", workspace, repoSlug, commit))
	flattenCommitDiffstat(&diffstatResponse, d)
	return nil
}

// CommitDiffstatResponse represents the response from the commit diffstat API
type CommitDiffstatResponse struct {
	Values []CommitDiffstat `json:"values"`
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Next   string           `json:"next"`
}

// CommitDiffstat represents diff statistics for a file in a commit
type CommitDiffstat struct {
	NewPath      string `json:"new_path"`
	OldPath      string `json:"old_path"`
	NewFile      bool   `json:"new_file"`
	RenamedFile  bool   `json:"renamed_file"`
	DeletedFile  bool   `json:"deleted_file"`
	LinesAdded   int    `json:"lines_added"`
	LinesRemoved int    `json:"lines_removed"`
	Type         string `json:"type"`
	Status       string `json:"status"`
}

// Flattens the commit diffstat information
func flattenCommitDiffstat(c *CommitDiffstatResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	diffstat := make([]interface{}, len(c.Values))
	for i, stat := range c.Values {
		diffstat[i] = map[string]interface{}{
			"new_path":      stat.NewPath,
			"old_path":      stat.OldPath,
			"new_file":      stat.NewFile,
			"renamed_file":  stat.RenamedFile,
			"deleted_file":  stat.DeletedFile,
			"lines_added":   stat.LinesAdded,
			"lines_removed": stat.LinesRemoved,
			"type":          stat.Type,
			"status":        stat.Status,
		}
	}

	d.Set("diffstat", diffstat)
}
