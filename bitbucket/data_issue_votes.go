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

func dataIssueVotes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueVotesRead,
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

func dataIssueVotesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	issueID := d.Get("issue_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueVotesRead", dumpResourceData(d, dataIssueVotes().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues/%s/votes", workspace, repoSlug, issueID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue votes call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue %s in repository %s/%s", issueID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue votes with params (%s): ", dumpResourceData(d, dataIssueVotes().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	votesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue votes response: %v", votesBody)

	var votesResponse IssueVotesResponse
	decodeerr := json.Unmarshal(votesBody, &votesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues/%s/votes", workspace, repoSlug, issueID))
	flattenIssueVotes(&votesResponse, d)
	return nil
}

// IssueVotesResponse represents the response from the issue votes API
type IssueVotesResponse struct {
	Values []IssueVote `json:"values"`
	Page   int         `json:"page"`
	Size   int         `json:"size"`
	Next   string      `json:"next"`
}

// IssueVote represents a vote on an issue
type IssueVote struct {
	User      map[string]interface{} `json:"user"`
	CreatedOn string                 `json:"created_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue votes information
func flattenIssueVotes(c *IssueVotesResponse, d *schema.ResourceData) {
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
