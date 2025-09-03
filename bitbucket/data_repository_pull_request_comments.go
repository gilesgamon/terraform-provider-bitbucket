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

func dataRepositoryPullRequestComments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryPullRequestCommentsRead,
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
			"comments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"content": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"user": {
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
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"parent": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"inline": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"links": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataRepositoryPullRequestCommentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pullRequestID := d.Get("pull_request_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryPullRequestCommentsRead", dumpResourceData(d, dataRepositoryPullRequestComments().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pullrequests/%s/comments", workspace, repoSlug, pullRequestID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository pull request comments call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pull request %s in repository %s/%s", pullRequestID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository pull request comments with params (%s): ", dumpResourceData(d, dataRepositoryPullRequestComments().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	commentsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository pull request comments response: %v", commentsBody)

	var commentsResponse RepositoryPullRequestCommentsResponse
	decodeerr := json.Unmarshal(commentsBody, &commentsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pullrequests/%s/comments", workspace, repoSlug, pullRequestID))
	flattenRepositoryPullRequestComments(&commentsResponse, d)
	return nil
}

// RepositoryPullRequestCommentsResponse represents the response from the repository pull request comments API
type RepositoryPullRequestCommentsResponse struct {
	Values []RepositoryPullRequestComment `json:"values"`
	Page   int                            `json:"page"`
	Size   int                            `json:"size"`
	Next   string                         `json:"next"`
}

// RepositoryPullRequestComment represents a pull request comment
type RepositoryPullRequestComment struct {
	ID        int                    `json:"id"`
	Content   map[string]interface{} `json:"content"`
	User      map[string]interface{} `json:"user"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Deleted   bool                   `json:"deleted"`
	Parent    map[string]interface{} `json:"parent"`
	Inline    map[string]interface{} `json:"inline"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the repository pull request comments information
func flattenRepositoryPullRequestComments(c *RepositoryPullRequestCommentsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	comments := make([]interface{}, len(c.Values))
	for i, comment := range c.Values {
		comments[i] = map[string]interface{}{
			"id":         comment.ID,
			"content":    comment.Content,
			"user":       comment.User,
			"created_on": comment.CreatedOn,
			"updated_on": comment.UpdatedOn,
			"deleted":    comment.Deleted,
			"parent":     comment.Parent,
			"inline":     comment.Inline,
			"links":      comment.Links,
		}
	}

	d.Set("comments", comments)
}
