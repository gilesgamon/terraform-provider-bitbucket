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
			"commit_sha": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Commit SHA to retrieve comments for",
			},
			"comments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_on": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inline": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"from": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"to": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"path": {
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

func dataCommitCommentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	commitSha := d.Get("commit_sha").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitCommentsRead", dumpResourceData(d, dataCommitComments().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commit/%s/comments",
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
		return diag.Errorf("no response returned from commit comments call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commitSha, workspace, repoSlug)
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

	d.SetId(fmt.Sprintf("%s/%s/%s/comments", workspace, repoSlug, commitSha))
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

// CommitComment represents a comment on a commit
type CommitComment struct {
	ID        string                 `json:"id"`
	Content   string                 `json:"content"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Author    *bitbucket.Account     `json:"author"`
	Type      string                 `json:"type"`
	Inline    *CommentInline         `json:"inline,omitempty"`
	Links     map[string]interface{} `json:"links"`
}

// CommentInline represents inline comment information
type CommentInline struct {
	From int    `json:"from"`
	To   int    `json:"to"`
	Path string `json:"path"`
}

// Flattens the commit comments information
func flattenCommitComments(c *CommitCommentsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	comments := make([]interface{}, len(c.Values))
	for i, comment := range c.Values {
		comments[i] = map[string]interface{}{
			"id":          comment.ID,
			"content":     comment.Content,
			"created_on":  comment.CreatedOn,
			"updated_on":  comment.UpdatedOn,
			"type":        comment.Type,
			"author":      flattenAccount(comment.Author),
			"inline":      flattenCommentInline(comment.Inline),
		}
	}

	d.Set("comments", comments)
}

// Flattens the comment inline information
func flattenCommentInline(inline *CommentInline) []interface{} {
	if inline == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"from": inline.From,
			"to":   inline.To,
			"path": inline.Path,
		},
	}
}
