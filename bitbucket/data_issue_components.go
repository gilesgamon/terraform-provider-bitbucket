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

func dataIssueComponents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueComponentsRead,
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

func dataIssueComponentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueComponentsRead", dumpResourceData(d, dataIssueComponents().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/components", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue components call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue components with params (%s): ", dumpResourceData(d, dataIssueComponents().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	componentsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue components response: %v", componentsBody)

	var componentsResponse IssueComponentsResponse
	decodeerr := json.Unmarshal(componentsBody, &componentsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/components", workspace, repoSlug))
	flattenIssueComponents(&componentsResponse, d)
	return nil
}

// IssueComponentsResponse represents the response from the issue components API
type IssueComponentsResponse struct {
	Values []IssueComponent `json:"values"`
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Next   string           `json:"next"`
}

// IssueComponent represents an issue component
type IssueComponent struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the issue components information
func flattenIssueComponents(c *IssueComponentsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	components := make([]interface{}, len(c.Values))
	for i, component := range c.Values {
		components[i] = map[string]interface{}{
			"uuid":        component.UUID,
			"name":        component.Name,
			"description": component.Description,
			"created_on":  component.CreatedOn,
			"updated_on":  component.UpdatedOn,
			"links":       component.Links,
		}
	}

	d.Set("components", components)
}
