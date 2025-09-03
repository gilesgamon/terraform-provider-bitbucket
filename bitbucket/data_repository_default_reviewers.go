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

func dataRepositoryDefaultReviewers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryDefaultReviewersRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_reviewers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nickname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_staff": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"account_status": {
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

func dataRepositoryDefaultReviewersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryDefaultReviewersRead", dumpResourceData(d, dataRepositoryDefaultReviewers().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/default-reviewers", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository default reviewers call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository default reviewers with params (%s): ", dumpResourceData(d, dataRepositoryDefaultReviewers().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	reviewersBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository default reviewers response: %v", reviewersBody)

	var reviewersResponse RepositoryDefaultReviewersResponse
	decodeerr := json.Unmarshal(reviewersBody, &reviewersResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/default-reviewers", workspace, repoSlug))
	flattenRepositoryDefaultReviewers(&reviewersResponse, d)
	return nil
}

// RepositoryDefaultReviewersResponse represents the response from the repository default reviewers API
type RepositoryDefaultReviewersResponse struct {
	Values []RepositoryDefaultReviewer `json:"values"`
	Page   int                         `json:"page"`
	Size   int                         `json:"size"`
	Next   string                      `json:"next"`
}

// RepositoryDefaultReviewer represents a default reviewer in a repository
type RepositoryDefaultReviewer struct {
	UUID          string                 `json:"uuid"`
	Username      string                 `json:"username"`
	DisplayName   string                 `json:"display_name"`
	Type          string                 `json:"type"`
	Nickname      string                 `json:"nickname"`
	AccountID     string                 `json:"account_id"`
	CreatedOn     string                 `json:"created_on"`
	IsStaff       bool                   `json:"is_staff"`
	AccountStatus string                 `json:"account_status"`
	Links         map[string]interface{} `json:"links"`
}

// Flattens the repository default reviewers information
func flattenRepositoryDefaultReviewers(c *RepositoryDefaultReviewersResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	reviewers := make([]interface{}, len(c.Values))
	for i, reviewer := range c.Values {
		reviewers[i] = map[string]interface{}{
			"uuid":           reviewer.UUID,
			"username":       reviewer.Username,
			"display_name":   reviewer.DisplayName,
			"type":           reviewer.Type,
			"nickname":       reviewer.Nickname,
			"account_id":     reviewer.AccountID,
			"created_on":     reviewer.CreatedOn,
			"is_staff":       reviewer.IsStaff,
			"account_status": reviewer.AccountStatus,
			"links":          reviewer.Links,
		}
	}

	d.Set("default_reviewers", reviewers)
}
