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

func dataRepositoryPullRequestActivity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryPullRequestActivityRead,
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
			"activity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
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
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"comment": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"changes_requested": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"approved": {
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

func dataRepositoryPullRequestActivityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pullRequestID := d.Get("pull_request_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryPullRequestActivityRead", dumpResourceData(d, dataRepositoryPullRequestActivity().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pullrequests/%s/activity", workspace, repoSlug, pullRequestID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository pull request activity call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pull request %s in repository %s/%s", pullRequestID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository pull request activity with params (%s): ", dumpResourceData(d, dataRepositoryPullRequestActivity().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	activityBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository pull request activity response: %v", activityBody)

	var activityResponse RepositoryPullRequestActivityResponse
	decodeerr := json.Unmarshal(activityBody, &activityResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pullrequests/%s/activity", workspace, repoSlug, pullRequestID))
	flattenRepositoryPullRequestActivity(&activityResponse, d)
	return nil
}

// RepositoryPullRequestActivityResponse represents the response from the repository pull request activity API
type RepositoryPullRequestActivityResponse struct {
	Values []RepositoryPullRequestActivityItem `json:"values"`
	Page   int                                 `json:"page"`
	Size   int                                 `json:"size"`
	Next   string                              `json:"next"`
}

// RepositoryPullRequestActivityItem represents a pull request activity item
type RepositoryPullRequestActivityItem struct {
	ID               int                    `json:"id"`
	Type             string                 `json:"type"`
	User             map[string]interface{} `json:"user"`
	CreatedOn        string                 `json:"created_on"`
	Update           map[string]interface{} `json:"update"`
	Comment          map[string]interface{} `json:"comment"`
	ChangesRequested map[string]interface{} `json:"changes_requested"`
	Approved         map[string]interface{} `json:"approved"`
	Links            map[string]interface{} `json:"links"`
}

// Flattens the repository pull request activity information
func flattenRepositoryPullRequestActivity(c *RepositoryPullRequestActivityResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	activity := make([]interface{}, len(c.Values))
	for i, item := range c.Values {
		activity[i] = map[string]interface{}{
			"id":                item.ID,
			"type":              item.Type,
			"user":              item.User,
			"created_on":        item.CreatedOn,
			"update":            item.Update,
			"comment":           item.Comment,
			"changes_requested": item.ChangesRequested,
			"approved":          item.Approved,
			"links":             item.Links,
		}
	}

	d.Set("activity", activity)
}
