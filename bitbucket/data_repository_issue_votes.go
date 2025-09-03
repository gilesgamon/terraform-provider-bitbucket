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

func dataRepositoryIssueVotes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryIssueVotesRead,
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Issue ID",
			},
			"votes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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

func dataRepositoryIssueVotesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	issueID := d.Get("issue_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryIssueVotesRead", dumpResourceData(d, dataRepositoryIssueVotes().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues/%s/votes", workspace, repoSlug, issueID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository issue votes call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue %s in repository %s/%s", issueID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository issue votes with params (%s): ", dumpResourceData(d, dataRepositoryIssueVotes().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	votesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository issue votes response: %v", votesBody)

	var votesResponse RepositoryIssueVotesResponse
	decodeerr := json.Unmarshal(votesBody, &votesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues/%s/votes", workspace, repoSlug, issueID))
	flattenRepositoryIssueVotes(&votesResponse, d)
	return nil
}

// RepositoryIssueVotesResponse represents the response from the repository issue votes API
type RepositoryIssueVotesResponse struct {
	Values []RepositoryIssueVote `json:"values"`
	Page   int                   `json:"page"`
	Size   int                   `json:"size"`
	Next   string                `json:"next"`
}

// RepositoryIssueVote represents an issue vote
type RepositoryIssueVote struct {
	User      map[string]interface{} `json:"user"`
	CreatedOn string                 `json:"created_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the repository issue votes information
func flattenRepositoryIssueVotes(c *RepositoryIssueVotesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	votes := make([]interface{}, len(c.Values))
	for i, vote := range c.Values {
		votes[i] = map[string]interface{}{
			"user":       vote.User,
			"created_on": vote.CreatedOn,
			"links":      vote.Links,
		}
	}

	d.Set("votes", votes)
}
