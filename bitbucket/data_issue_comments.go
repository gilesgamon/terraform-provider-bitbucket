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

func dataIssueComments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueCommentsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issue_id": {
				Type:     schema.TypeString,
				Required: true,
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

func dataIssueCommentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	issueID := d.Get("issue_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueCommentsRead", dumpResourceData(d, dataIssueComments().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues/%s/comments", workspace, repoSlug, issueID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue comments call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue %s in repository %s/%s", issueID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue comments with params (%s): ", dumpResourceData(d, dataIssueComments().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	commentsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue comments response: %v", commentsBody)

	var commentsResponse IssueCommentsResponse
	decodeerr := json.Unmarshal(commentsBody, &commentsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues/%s/comments", workspace, repoSlug, issueID))
	flattenIssueComments(&commentsResponse, d)
	return nil
}

// IssueCommentsResponse represents the response from the issue comments API
type IssueCommentsResponse struct {
	Values []IssueComment `json:"values"`
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Next   string         `json:"next"`
}

// IssueComment represents a comment on an issue
type IssueComment struct {
	ID        int                    `json:"id"`
	Content   map[string]interface{} `json:"content"`
	User      map[string]interface{} `json:"user"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue comments information
func flattenIssueComments(c *IssueCommentsResponse, d *schema.ResourceData) {
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
			"links":      comment.Links,
		}
	}

	d.Set("comments", comments)
}
