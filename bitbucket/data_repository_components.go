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

func dataRepositoryComponents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryComponentsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"components": {
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
						"assignee": {
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

func dataRepositoryComponentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryComponentsRead", dumpResourceData(d, dataRepositoryComponents().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/components", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository components call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository components with params (%s): ", dumpResourceData(d, dataRepositoryComponents().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	componentsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository components response: %v", componentsBody)

	var componentsResponse RepositoryComponentsResponse
	decodeerr := json.Unmarshal(componentsBody, &componentsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/components", workspace, repoSlug))
	flattenRepositoryComponents(&componentsResponse, d)
	return nil
}

// RepositoryComponentsResponse represents the response from the repository components API
type RepositoryComponentsResponse struct {
	Values []RepositoryComponent `json:"values"`
	Page   int                   `json:"page"`
	Size   int                   `json:"size"`
	Next   string                `json:"next"`
}

// RepositoryComponent represents a component in a repository
type RepositoryComponent struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Assignee    map[string]interface{} `json:"assignee"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the repository components information
func flattenRepositoryComponents(c *RepositoryComponentsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	components := make([]interface{}, len(c.Values))
	for i, component := range c.Values {
		components[i] = map[string]interface{}{
			"id":          component.ID,
			"name":        component.Name,
			"description": component.Description,
			"assignee":    component.Assignee,
			"created_on":  component.CreatedOn,
			"updated_on":  component.UpdatedOn,
			"links":       component.Links,
		}
	}

	d.Set("components", components)
}
