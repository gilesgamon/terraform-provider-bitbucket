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

func dataCommitComments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCommitCommentsRead,
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
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user": {
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

func dataCommitCommentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	commit := d.Get("commit").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitCommentsRead", dumpResourceData(d, dataCommitComments().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commits/%s/comments", workspace, repoSlug, commit)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from commit comments call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commit, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading commit comments with params (%s): ", dumpResourceData(d, dataCommitComments().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	commentsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] commit comments response: %v", commentsBody)

	var commentsResponse CommitCommentsResponse
	decodeerr := json.Unmarshal(commentsBody, &commentsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/commits/%s/comments", workspace, repoSlug, commit))
	flattenCommitComments(&commentsResponse, d)
	return nil
}

// CommitCommentsResponse represents the response from the commit comments API
type CommitCommentsResponse struct {
	Values []CommitComment `json:"values"`
	Page   int             `json:"page"`
	Size   int             `json:"size"`
	Next   string          `json:"next"`
}

// CommitComment represents a commit comment
type CommitComment struct {
	ID        int                    `json:"id"`
	Content   map[string]interface{} `json:"content"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	User      map[string]interface{} `json:"user"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the commit comments information
func flattenCommitComments(c *CommitCommentsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	comments := make([]interface{}, len(c.Values))
	for i, comment := range c.Values {
		comments[i] = map[string]interface{}{
			"id":         comment.ID,
			"content":    comment.Content,
			"created_on": comment.CreatedOn,
			"updated_on": comment.UpdatedOn,
			"user":       comment.User,
			"links":      comment.Links,
		}
	}

	d.Set("comments", comments)
}
