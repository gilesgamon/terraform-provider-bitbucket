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

func dataRepositoryPullRequestDiff() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryPullRequestDiffRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pull_request_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pull request ID",
			},
			"context": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of context lines to show around changes",
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to specific file to get diff for",
			},
			"diff": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"new": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"old": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"hunks": {
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
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"line_old": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"line_new": {
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
				},
			},
		},
	}
}

func dataRepositoryPullRequestDiffRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pullRequestID := d.Get("pull_request_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryPullRequestDiffRead", dumpResourceData(d, dataRepositoryPullRequestDiff().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pullrequests/%s/diff", workspace, repoSlug, pullRequestID)

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
		return diag.Errorf("no response returned from repository pull request diff call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pull request %s in repository %s/%s", pullRequestID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository pull request diff with params (%s): ", dumpResourceData(d, dataRepositoryPullRequestDiff().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	diffBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository pull request diff response: %v", diffBody)

	var diffResponse RepositoryPullRequestDiffResponse
	decodeerr := json.Unmarshal(diffBody, &diffResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pullrequests/%s/diff", workspace, repoSlug, pullRequestID))
	flattenRepositoryPullRequestDiff(&diffResponse, d)
	return nil
}

// RepositoryPullRequestDiffResponse represents the response from the repository pull request diff API
type RepositoryPullRequestDiffResponse struct {
	Values []RepositoryPullRequestDiff `json:"values"`
	Page   int                         `json:"page"`
	Size   int                         `json:"size"`
	Next   string                      `json:"next"`
}

// RepositoryPullRequestDiff represents a pull request diff
type RepositoryPullRequestDiff struct {
	New   map[string]interface{} `json:"new"`
	Old   map[string]interface{} `json:"old"`
	Hunks []PRDiffHunk             `json:"hunks"`
}

// PRDiffHunk represents a diff hunk for pull requests
type PRDiffHunk struct {
	Type  string      `json:"type"`
	Lines []PRDiffLine `json:"lines"`
}

// PRDiffLine represents a diff line for pull requests
type PRDiffLine struct {
	LineOld int    `json:"line_old"`
	LineNew int    `json:"line_new"`
	Content string `json:"content"`
}

// Flattens the repository pull request diff information
func flattenRepositoryPullRequestDiff(c *RepositoryPullRequestDiffResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	diff := make([]interface{}, len(c.Values))
	for i, diffItem := range c.Values {
		hunks := make([]interface{}, len(diffItem.Hunks))
		for j, hunk := range diffItem.Hunks {
			lines := make([]interface{}, len(hunk.Lines))
			for k, line := range hunk.Lines {
				lines[k] = map[string]interface{}{
					"line_old": line.LineOld,
					"line_new": line.LineNew,
					"content":  line.Content,
				}
			}
			hunks[j] = map[string]interface{}{
				"type":  hunk.Type,
				"lines": lines,
			}
		}

		diff[i] = map[string]interface{}{
			"new":   diffItem.New,
			"old":   diffItem.Old,
			"hunks": hunks,
		}
	}

	d.Set("diff", diff)
}
