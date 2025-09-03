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

func dataCommitDiff() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCommitDiffRead,
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
			"diff": {
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
						"similarity": {
							Type:     schema.TypeInt,
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
						"hunks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"old_start": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"old_lines": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"new_start": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"new_lines": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"content": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataCommitDiffRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	commit := d.Get("commit").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitDiffRead", dumpResourceData(d, dataCommitDiff().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commits/%s/diff", workspace, repoSlug, commit)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from commit diff call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commit, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading commit diff with params (%s): ", dumpResourceData(d, dataCommitDiff().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	diffBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] commit diff response: %v", diffBody)

	var diffResponse CommitDiffResponse
	decodeerr := json.Unmarshal(diffBody, &diffResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/commits/%s/diff", workspace, repoSlug, commit))
	flattenCommitDiff(&diffResponse, d)
	return nil
}

// CommitDiffResponse represents the response from the commit diff API
type CommitDiffResponse struct {
	Values []CommitDiff `json:"values"`
	Page   int          `json:"page"`
	Size   int          `json:"size"`
	Next   string       `json:"next"`
}

// CommitDiff represents a file diff in a commit
type CommitDiff struct {
	NewPath      string       `json:"new_path"`
	OldPath      string       `json:"old_path"`
	NewFile      bool         `json:"new_file"`
	RenamedFile  bool         `json:"renamed_file"`
	DeletedFile  bool         `json:"deleted_file"`
	Similarity   int          `json:"similarity"`
	Status       string       `json:"status"`
	LinesAdded   int          `json:"lines_added"`
	LinesRemoved int          `json:"lines_removed"`
	Hunks        []DiffHunk   `json:"hunks"`
}

// DiffHunk represents a hunk of changes in a diff
type DiffHunk struct {
	OldStart int    `json:"old_start"`
	OldLines int    `json:"old_lines"`
	NewStart int    `json:"new_start"`
	NewLines int    `json:"new_lines"`
	Content  string `json:"content"`
}

// Flattens the commit diff information
func flattenCommitDiff(c *CommitDiffResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	diff := make([]interface{}, len(c.Values))
	for i, fileDiff := range c.Values {
		hunks := make([]interface{}, len(fileDiff.Hunks))
		for j, hunk := range fileDiff.Hunks {
			hunks[j] = map[string]interface{}{
				"old_start": hunk.OldStart,
				"old_lines": hunk.OldLines,
				"new_start": hunk.NewStart,
				"new_lines": hunk.NewLines,
				"content":   hunk.Content,
			}
		}

		diff[i] = map[string]interface{}{
			"new_path":      fileDiff.NewPath,
			"old_path":      fileDiff.OldPath,
			"new_file":      fileDiff.NewFile,
			"renamed_file":  fileDiff.RenamedFile,
			"deleted_file":  fileDiff.DeletedFile,
			"similarity":    fileDiff.Similarity,
			"status":        fileDiff.Status,
			"lines_added":   fileDiff.LinesAdded,
			"lines_removed": fileDiff.LinesRemoved,
			"hunks":         hunks,
		}
	}

	d.Set("diff", diff)
}
