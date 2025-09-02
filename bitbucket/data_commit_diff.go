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
			"commit_sha": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Commit SHA to retrieve diff for",
			},
			"context": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of context lines to show around changes",
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to filter diff by",
			},
			"diff": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"new_file": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"deleted_file": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"renamed_file": {
							Type:     schema.TypeBool,
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
									"context": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"segments": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"lines": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"truncated": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"stats": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
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
	commitSha := d.Get("commit_sha").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitDiffRead", dumpResourceData(d, dataCommitDiff().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/diff/%s", workspace, repoSlug, commitSha)

	// Build query parameters
	params := make(map[string]string)
	if context, ok := d.GetOk("context"); ok {
		params["context"] = fmt.Sprintf("%d", context.(int))
	}
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
		return diag.Errorf("no response returned from commit diff call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commitSha, workspace, repoSlug)
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

	d.SetId(fmt.Sprintf("%s/%s/%s/diff", workspace, repoSlug, commitSha))
	flattenCommitDiff(&diffResponse, d)
	return nil
}

// CommitDiffResponse represents the response from the commit diff API
type CommitDiffResponse struct {
	Values []DiffFile `json:"values"`
	Page   int        `json:"page"`
	Size   int        `json:"size"`
	Next   string     `json:"next"`
}

// DiffFile represents a file in the commit diff
type DiffFile struct {
	NewFile     bool       `json:"new_file"`
	DeletedFile bool       `json:"deleted_file"`
	RenamedFile bool       `json:"renamed_file"`
	OldPath     string     `json:"old_path"`
	NewPath     string     `json:"new_path"`
	Hunks       []DiffHunk `json:"hunks"`
	Stats       map[string]interface{} `json:"stats"`
	Links       map[string]interface{} `json:"links"`
}

// DiffHunk represents a hunk of changes in a file
type DiffHunk struct {
	OldStart int    `json:"old_start"`
	OldLines int    `json:"old_lines"`
	NewStart int    `json:"new_start"`
	NewLines int    `json:"new_lines"`
	Context  string `json:"context"`
	Segments []DiffSegment `json:"segments"`
}

// DiffSegment represents a segment of lines in a diff hunk
type DiffSegment struct {
	Type      string   `json:"type"`
	Lines     []string `json:"lines"`
	Truncated bool     `json:"truncated"`
}

// Flattens the commit diff information
func flattenCommitDiff(c *CommitDiffResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	diff := make([]interface{}, len(c.Values))
	for i, file := range c.Values {
		hunks := make([]interface{}, len(file.Hunks))
		for j, hunk := range file.Hunks {
			segments := make([]interface{}, len(hunk.Segments))
			for k, segment := range hunk.Segments {
				segments[k] = map[string]interface{}{
					"type":      segment.Type,
					"lines":     segment.Lines,
					"truncated": segment.Truncated,
				}
			}
			hunks[j] = map[string]interface{}{
				"old_start": hunk.OldStart,
				"old_lines": hunk.OldLines,
				"new_start": hunk.NewStart,
				"new_lines": hunk.NewLines,
				"context":   hunk.Context,
				"segments":  segments,
			}
		}

		diff[i] = map[string]interface{}{
			"new_file":     file.NewFile,
			"deleted_file": file.DeletedFile,
			"renamed_file": file.RenamedFile,
			"old_path":     file.OldPath,
			"new_path":     file.NewPath,
			"hunks":        hunks,
			"stats":        file.Stats,
		}
	}

	d.Set("diff", diff)
}
