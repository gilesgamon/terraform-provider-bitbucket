package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/DrFaust92/bitbucket-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataPullRequest() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPullRequestRead,
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
				Description: "Pull request ID (number)",
			},
			"id": {
				Type:     schema.TypeString,
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
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"source": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"branch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"destination": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"branch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"merge_commit": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataPullRequestRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pullRequestID := d.Get("pull_request_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPullRequestRead", dumpResourceData(d, dataPullRequest().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pullrequests/%s",
		workspace,
		repoSlug,
		pullRequestID,
	)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pull request call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pull request %s in repository %s/%s", pullRequestID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pull request information with params (%s): ", dumpResourceData(d, dataPullRequest().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	prBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pull request response: %v", prBody)

	var pr PullRequest
	decodeerr := json.Unmarshal(prBody, &pr)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/%d", workspace, repoSlug, pr.ID))
	flattenPullRequest(&pr, d)
	return nil
}

// PullRequest represents a Bitbucket pull request
type PullRequest struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	State       string                 `json:"state"`
	Author      *bitbucket.Account     `json:"author"`
	Source      PullRequestBranch      `json:"source"`
	Destination PullRequestBranch      `json:"destination"`
	Reviewers   []PullRequestReviewer  `json:"reviewers"`
	Links       map[string]interface{} `json:"links"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	ClosedBy    *bitbucket.Account     `json:"closed_by,omitempty"`
	MergedBy    *bitbucket.Account     `json:"merged_by,omitempty"`
	MergeCommit *PullRequestCommit     `json:"merge_commit,omitempty"`
}

// PullRequestBranch represents a branch in a pull request
type PullRequestBranch struct {
	Branch     PullRequestBranchInfo `json:"branch"`
	Commit     PullRequestCommit     `json:"commit"`
	Repository bitbucket.Repository  `json:"repository"`
}

// PullRequestBranchInfo represents branch information
type PullRequestBranchInfo struct {
	Name string `json:"name"`
}

// PullRequestCommit represents a commit in a pull request
type PullRequestCommit struct {
	Hash  string                 `json:"hash"`
	Type  string                 `json:"type"`
	Links map[string]interface{} `json:"links"`
}

// PullRequestReviewer represents a reviewer in a pull request
type PullRequestReviewer struct {
	User     bitbucket.Account `json:"user"`
	Type     string            `json:"type"`
	Approved bool               `json:"approved"`
}

// Flattens the pull request information
func flattenPullRequest(pr *PullRequest, d *schema.ResourceData) {
	if pr == nil {
		return
	}

	d.Set("id", fmt.Sprintf("%d", pr.ID))
	d.Set("title", pr.Title)
	d.Set("description", pr.Description)
	d.Set("state", pr.State)
	d.Set("created_date", pr.CreatedOn)
	d.Set("updated_date", pr.UpdatedOn)
	d.Set("author", flattenPullRequestAccount(pr.Author))
	d.Set("source", flattenPullRequestBranch(pr.Source))
	d.Set("destination", flattenPullRequestBranch(pr.Destination))
}

// Flattens the pull request branch information
func flattenPullRequestBranch(b PullRequestBranch) []interface{} {
	if b.Branch.Name == "" {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"branch":     b.Branch.Name,
			"commit":     b.Commit.Hash,
			"repository": b.Repository.Name,
		},
	}
}

// Flattens a pull request account (for author, closed_by, merged_by)
func flattenPullRequestAccount(a *bitbucket.Account) []interface{} {
	if a == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"username":     a.Username,
			"display_name": a.DisplayName,
			"uuid":         a.Uuid,
		},
	}
}
