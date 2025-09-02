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
			"branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Branch to get commits from (defaults to main/master)",
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to filter commits by",
			},
			"include": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Include additional data (e.g., 'include=stats')",
			},
			"exclude": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Exclude certain data (e.g., 'exclude=stats')",
			},
			"merges": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter merge commits ('include', 'exclude', or 'only')",
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"author_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"author_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"committer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"committer_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"committer_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"summary": {
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
						"date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rendered": {
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

	// Build query parameters
	params := make(map[string]string)
	if branch, ok := d.GetOk("branch"); ok {
		params["include"] = branch.(string)
	}
	if path, ok := d.GetOk("path"); ok {
		params["path"] = path.(string)
	}
	if include, ok := d.GetOk("include"); ok {
		params["include"] = include.(string)
	}
	if exclude, ok := d.GetOk("exclude"); ok {
		params["exclude"] = exclude.(string)
	}
	if merges, ok := d.GetOk("merges"); ok {
		params["merges"] = merges.(string)
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

// CommitListItem represents a commit in the repository list
type CommitListItem struct {
	Hash           string                 `json:"hash"`
	Author         string                 `json:"author"`
	AuthorEmail    string                 `json:"author_email"`
	AuthorDate     string                 `json:"author_date"`
	Committer      string                 `json:"committer"`
	CommitterEmail string                 `json:"committer_email"`
	CommitterDate  string                 `json:"committer_date"`
	Message        string                 `json:"message"`
	Summary        string                 `json:"summary"`
	Parents        []string               `json:"parents"`
	Date           string                 `json:"date"`
	Rendered       map[string]interface{} `json:"rendered"`
	Links          map[string]interface{} `json:"links"`
}

// Flattens the commits information
func flattenCommits(c *CommitsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	commits := make([]interface{}, len(c.Values))
	for i, commit := range c.Values {
		commits[i] = map[string]interface{}{
			"hash":            commit.Hash,
			"author":          commit.Author,
			"author_email":    commit.AuthorEmail,
			"author_date":     commit.AuthorDate,
			"committer":       commit.Committer,
			"committer_email": commit.CommitterEmail,
			"committer_date":  commit.CommitterDate,
			"message":         commit.Message,
			"summary":         commit.Summary,
			"parents":         commit.Parents,
			"date":            commit.Date,
			"rendered":        commit.Rendered,
		}
	}

	d.Set("commits", commits)
}
