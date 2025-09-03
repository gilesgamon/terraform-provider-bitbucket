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

func dataCommits() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCommitsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commits": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"author": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parents": {
							Type:     schema.TypeList,
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

func dataCommitsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitsRead", dumpResourceData(d, dataCommits().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commits", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from commits call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading commits with params (%s): ", dumpResourceData(d, dataCommits().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	commitsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] commits response: %v", commitsBody)

	var commitsResponse CommitsResponse
	decodeerr := json.Unmarshal(commitsBody, &commitsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/commits", workspace, repoSlug))
	flattenCommits(&commitsResponse, d)
	return nil
}

// CommitsResponse represents the response from the commits API
type CommitsResponse struct {
	Values []CommitListItem `json:"values"`
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Next   string           `json:"next"`
}

// CommitListItem represents a commit in the list
type CommitListItem struct {
	Hash    string                 `json:"hash"`
	Author  map[string]interface{} `json:"author"`
	Message string                 `json:"message"`
	Date    string                 `json:"date"`
	Parents []string               `json:"parents"`
	Links   map[string]interface{} `json:"links"`
}

// Flattens the commits information
func flattenCommits(c *CommitsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	commits := make([]interface{}, len(c.Values))
	for i, commit := range c.Values {
		commits[i] = map[string]interface{}{
			"hash":    commit.Hash,
			"author":  commit.Author,
			"message": commit.Message,
			"date":    commit.Date,
			"parents": commit.Parents,
			"links":   commit.Links,
		}
	}

	d.Set("commits", commits)
}
