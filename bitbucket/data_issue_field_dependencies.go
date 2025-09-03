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

func dataIssueFieldDependencies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldDependenciesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"field_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dependencies": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"field": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"value": {
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

func dataIssueFieldDependenciesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldDependenciesRead", dumpResourceData(d, dataIssueFieldDependencies().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/dependencies", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field dependencies call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field dependencies with params (%s): ", dumpResourceData(d, dataIssueFieldDependencies().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	dependenciesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field dependencies response: %v", dependenciesBody)

	var dependenciesResponse IssueFieldDependenciesResponse
	decodeerr := json.Unmarshal(dependenciesBody, &dependenciesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/dependencies", workspace, repoSlug, fieldUUID))
	flattenIssueFieldDependencies(&dependenciesResponse, d)
	return nil
}

// IssueFieldDependenciesResponse represents the response from the issue field dependencies API
type IssueFieldDependenciesResponse struct {
	Values []IssueFieldDependency `json:"values"`
	Page   int                    `json:"page"`
	Size   int                    `json:"size"`
	Next   string                 `json:"next"`
}

// IssueFieldDependency represents a dependency for an issue field
type IssueFieldDependency struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Field     map[string]interface{} `json:"field"`
	Value     string                 `json:"value"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field dependencies information
func flattenIssueFieldDependencies(c *IssueFieldDependenciesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	dependencies := make([]interface{}, len(c.Values))
	for i, dependency := range c.Values {
		dependencies[i] = map[string]interface{}{
			"uuid":       dependency.UUID,
			"name":       dependency.Name,
			"type":       dependency.Type,
			"field":      dependency.Field,
			"value":      dependency.Value,
			"created_on": dependency.CreatedOn,
			"updated_on": dependency.UpdatedOn,
			"links":      dependency.Links,
		}
	}

	d.Set("dependencies", dependencies)
}
