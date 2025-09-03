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

func dataRepositoryMilestones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryMilestonesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"milestones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"release_date": {
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

func dataRepositoryMilestonesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryMilestonesRead", dumpResourceData(d, dataRepositoryMilestones().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/milestones", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository milestones call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository milestones with params (%s): ", dumpResourceData(d, dataRepositoryMilestones().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	milestonesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository milestones response: %v", milestonesBody)

	var milestonesResponse RepositoryMilestonesResponse
	decodeerr := json.Unmarshal(milestonesBody, &milestonesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/milestones", workspace, repoSlug))
	flattenRepositoryMilestones(&milestonesResponse, d)
	return nil
}

// RepositoryMilestonesResponse represents the response from the repository milestones API
type RepositoryMilestonesResponse struct {
	Values []RepositoryMilestone `json:"values"`
	Page   int                   `json:"page"`
	Size   int                   `json:"size"`
	Next   string                `json:"next"`
}

// RepositoryMilestone represents a milestone in a repository
type RepositoryMilestone struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	State       string                 `json:"state"`
	StartDate   string                 `json:"start_date"`
	ReleaseDate string                 `json:"release_date"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the repository milestones information
func flattenRepositoryMilestones(c *RepositoryMilestonesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	milestones := make([]interface{}, len(c.Values))
	for i, milestone := range c.Values {
		milestones[i] = map[string]interface{}{
			"id":           milestone.ID,
			"name":         milestone.Name,
			"description":  milestone.Description,
			"state":        milestone.State,
			"start_date":   milestone.StartDate,
			"release_date": milestone.ReleaseDate,
			"created_on":   milestone.CreatedOn,
			"updated_on":   milestone.UpdatedOn,
			"links":        milestone.Links,
		}
	}

	d.Set("milestones", milestones)
}
