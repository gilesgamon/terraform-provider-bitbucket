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

func dataPullRequests() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPullRequestsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter PRs by state (open, merged, declined, superseded)",
			},
			"source_branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter PRs by source branch name",
			},
			"destination_branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter PRs by destination branch name",
			},
			"author": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter PRs by author username",
			},
			"reviewer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter PRs by reviewer username",
			},
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search query string",
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort field (created_on, updated_on, title, author)",
			},
			"pull_requests": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"author": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"destination": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"merge_commit": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"closed_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"closed_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reviewers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"display_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"approved": {
										Type:     schema.TypeBool,
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

func dataPullRequestsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPullRequestsRead", dumpResourceData(d, dataPullRequests().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pullrequests", workspace, repoSlug)

	// Build query parameters
	params := make(map[string]string)
	if state, ok := d.GetOk("state"); ok {
		params["state"] = state.(string)
	}
	if sourceBranch, ok := d.GetOk("source_branch"); ok {
		params["source_branch"] = sourceBranch.(string)
	}
	if destinationBranch, ok := d.GetOk("destination_branch"); ok {
		params["destination_branch"] = destinationBranch.(string)
	}
	if author, ok := d.GetOk("author"); ok {
		params["author"] = author.(string)
	}
	if reviewer, ok := d.GetOk("reviewer"); ok {
		params["reviewer"] = reviewer.(string)
	}
	if q, ok := d.GetOk("q"); ok {
		params["q"] = q.(string)
	}
	if sort, ok := d.GetOk("sort"); ok {
		params["sort"] = sort.(string)
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
		return diag.Errorf("no response returned from pull requests call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pull requests with params (%s): ", dumpResourceData(d, dataPullRequests().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	pullRequestsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pull requests response: %v", pullRequestsBody)

	var pullRequestsResponse PullRequestsResponse
	decodeerr := json.Unmarshal(pullRequestsBody, &pullRequestsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pullrequests", workspace, repoSlug))
	flattenPullRequests(&pullRequestsResponse, d)
	return nil
}

// PullRequestsResponse represents the response from the pull requests API
type PullRequestsResponse struct {
	Values []PullRequestListItem `json:"values"`
	Page   int                   `json:"page"`
	Size   int                   `json:"size"`
	Next   string                `json:"next"`
}

// PullRequestListItem represents a pull request in the list
type PullRequestListItem struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	State       string                 `json:"state"`
	Author      string                 `json:"author"`
	Source      map[string]interface{} `json:"source"`
	Destination map[string]interface{} `json:"destination"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	MergeCommit map[string]interface{} `json:"merge_commit"`
	ClosedBy    string                 `json:"closed_by"`
	ClosedOn    string                 `json:"closed_on"`
	Reviewers   []PullRequestListItemReviewer  `json:"reviewers"`
	Links       map[string]interface{} `json:"links"`
}

// PullRequestListItemReviewer represents a reviewer on a pull request list item
type PullRequestListItemReviewer struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	Approved    bool   `json:"approved"`
}

// Flattens the pull requests information
func flattenPullRequests(c *PullRequestsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	pullRequests := make([]interface{}, len(c.Values))
	for i, pr := range c.Values {
		reviewers := make([]interface{}, len(pr.Reviewers))
		for j, reviewer := range pr.Reviewers {
			reviewers[j] = map[string]interface{}{
				"username":     reviewer.Username,
				"display_name": reviewer.DisplayName,
				"type":         reviewer.Type,
				"approved":     reviewer.Approved,
			}
		}

		pullRequests[i] = map[string]interface{}{
			"id":           pr.ID,
			"title":        pr.Title,
			"description":  pr.Description,
			"state":        pr.State,
			"author":       pr.Author,
			"source":       pr.Source,
			"destination":  pr.Destination,
			"created_on":   pr.CreatedOn,
			"updated_on":   pr.UpdatedOn,
			"merge_commit": pr.MergeCommit,
			"closed_by":    pr.ClosedBy,
			"closed_on":    pr.ClosedOn,
			"reviewers":    reviewers,
		}
	}

	d.Set("pull_requests", pullRequests)
}
