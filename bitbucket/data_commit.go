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

func dataCommit() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCommitRead,
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
				Description: "Commit SHA or branch name (e.g., 'main', 'develop', 'abc123...')",
			},
			"hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date": {
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
			"parents": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataCommitRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	commitSha := d.Get("commit_sha").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitRead", dumpResourceData(d, dataCommit().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commit/%s",
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
		return diag.Errorf("no response returned from commit call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commitSha, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading commit information with params (%s): ", dumpResourceData(d, dataCommit().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	commitBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] commit response: %v", commitBody)

	var commit Commit
	decodeerr := json.Unmarshal(commitBody, &commit)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", workspace, repoSlug, commit.Hash))
	flattenCommit(&commit, d)
	return nil
}

// Commit represents a Bitbucket commit
type Commit struct {
	Hash       string                 `json:"hash"`
	Type       string                 `json:"type"`
	Message    string                 `json:"message"`
	Author     CommitAuthor           `json:"author"`
	Committer  CommitAuthor           `json:"committer"`
	Date       string                 `json:"date"`
	Parents    []CommitParent         `json:"parents"`
	Repository bitbucket.Repository   `json:"repository"`
	Links      map[string]interface{} `json:"links"`
}

// CommitAuthor represents the author or committer of a commit
type CommitAuthor struct {
	Raw  string            `json:"raw"`
	Type string            `json:"type"`
	User bitbucket.Account `json:"user"`
}

// CommitParent represents a parent commit
type CommitParent struct {
	Hash  string                 `json:"hash"`
	Type  string                 `json:"type"`
	Links map[string]interface{} `json:"links"`
}

// Flattens the commit information
func flattenCommit(c *Commit, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("hash", c.Hash)
	d.Set("message", c.Message)
	d.Set("date", c.Date)
	d.Set("author", flattenCommitAuthor(c.Author))
	d.Set("parents", flattenCommitParents(c.Parents))
}

// Flattens the commit author information
func flattenCommitAuthor(a CommitAuthor) []interface{} {
	if a.Raw == "" {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"username":     a.User.Username,
			"display_name": a.User.DisplayName,
			"uuid":         a.User.Uuid,
		},
	}
}

// Flattens the commit parents
func flattenCommitParents(parents []CommitParent) []interface{} {
	if len(parents) == 0 {
		return nil
	}
	result := make([]interface{}, len(parents))
	for i, parent := range parents {
		result[i] = map[string]interface{}{
			"hash": parent.Hash,
			"type": parent.Type,
		}
	}
	return result
}
