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

func dataCommitPullRequests() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCommitPullRequestsRead,
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
				Description: "Commit SHA to find pull requests for",
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
					},
				},
			},
		},
	}
}

func dataCommitPullRequestsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	commitSha := d.Get("commit_sha").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitPullRequestsRead", dumpResourceData(d, dataCommitPullRequests().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commit/%s/pullrequests",
		workspace,
		repoSlug,
		commitSha,
	)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from commit pull requests call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commitSha, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading commit pull requests with params (%s): ", dumpResourceData(d, dataCommitPullRequests().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	pullRequestsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] commit pull requests response: %v", pullRequestsBody)

	var pullRequestsResponse CommitPullRequestsResponse
	decodeerr := json.Unmarshal(pullRequestsBody, &pullRequestsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/pullrequests", workspace, repoSlug, commitSha))
	flattenCommitPullRequests(&pullRequestsResponse, d)
	return nil
}

// CommitPullRequestsResponse represents the response from the commit pull requests API
type CommitPullRequestsResponse struct {
	Values []CommitPullRequest `json:"values"`
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	Next   string              `json:"next"`
}

// CommitPullRequest represents a pull request containing a specific commit
type CommitPullRequest struct {
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
	Links       map[string]interface{} `json:"links"`
}

// Flattens the commit pull requests information
func flattenCommitPullRequests(c *CommitPullRequestsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	pullRequests := make([]interface{}, len(c.Values))
	for i, pr := range c.Values {
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
		}
	}

	d.Set("pull_requests", pullRequests)
}
