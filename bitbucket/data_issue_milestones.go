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

func dataIssueMilestones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueMilestonesRead,
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
						"uuid": {
							Type:     schema.TypeString,
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

func dataIssueMilestonesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueMilestonesRead", dumpResourceData(d, dataIssueMilestones().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/milestones", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue milestones call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue milestones with params (%s): ", dumpResourceData(d, dataIssueMilestones().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	milestonesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue milestones response: %v", milestonesBody)

	var milestonesResponse IssueMilestonesResponse
	decodeerr := json.Unmarshal(milestonesBody, &milestonesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/milestones", workspace, repoSlug))
	flattenIssueMilestones(&milestonesResponse, d)
	return nil
}

// IssueMilestonesResponse represents the response from the issue milestones API
type IssueMilestonesResponse struct {
	Values []IssueMilestone `json:"values"`
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Next   string           `json:"next"`
}

// IssueMilestone represents an issue milestone
type IssueMilestone struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	State       string                 `json:"state"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the issue milestones information
func flattenIssueMilestones(c *IssueMilestonesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	milestones := make([]interface{}, len(c.Values))
	for i, milestone := range c.Values {
		milestones[i] = map[string]interface{}{
			"uuid":        milestone.UUID,
			"name":        milestone.Name,
			"description": milestone.Description,
			"state":       milestone.State,
			"created_on":  milestone.CreatedOn,
			"updated_on":  milestone.UpdatedOn,
			"links":       milestone.Links,
		}
	}

	d.Set("milestones", milestones)
}
