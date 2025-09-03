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

func dataRepositoryDeploymentEnvironments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryDeploymentEnvironmentsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environments": {
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
						"environment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rank": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"deployment_gate_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"deployment_gate_check": {
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

func dataRepositoryDeploymentEnvironmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryDeploymentEnvironmentsRead", dumpResourceData(d, dataRepositoryDeploymentEnvironments().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/environments", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository deployment environments call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository deployment environments with params (%s): ", dumpResourceData(d, dataRepositoryDeploymentEnvironments().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	environmentsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository deployment environments response: %v", environmentsBody)

	var environmentsResponse RepositoryDeploymentEnvironmentsResponse
	decodeerr := json.Unmarshal(environmentsBody, &environmentsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/environments", workspace, repoSlug))
	flattenRepositoryDeploymentEnvironments(&environmentsResponse, d)
	return nil
}

// RepositoryDeploymentEnvironmentsResponse represents the response from the repository deployment environments API
type RepositoryDeploymentEnvironmentsResponse struct {
	Values []RepositoryDeploymentEnvironment `json:"values"`
	Page   int                               `json:"page"`
	Size   int                               `json:"size"`
	Next   string                            `json:"next"`
}

// RepositoryDeploymentEnvironment represents a deployment environment
type RepositoryDeploymentEnvironment struct {
	UUID                    string                 `json:"uuid"`
	Name                    string                 `json:"name"`
	EnvironmentType         string                 `json:"environment_type"`
	Rank                    int                    `json:"rank"`
	DeploymentGateEnabled   bool                   `json:"deployment_gate_enabled"`
	DeploymentGateCheck     map[string]interface{} `json:"deployment_gate_check"`
	CreatedOn               string                 `json:"created_on"`
	UpdatedOn               string                 `json:"updated_on"`
	Links                   map[string]interface{} `json:"links"`
}

// Flattens the repository deployment environments information
func flattenRepositoryDeploymentEnvironments(c *RepositoryDeploymentEnvironmentsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	environments := make([]interface{}, len(c.Values))
	for i, environment := range c.Values {
		environments[i] = map[string]interface{}{
			"uuid":                    environment.UUID,
			"name":                    environment.Name,
			"environment_type":        environment.EnvironmentType,
			"rank":                    environment.Rank,
			"deployment_gate_enabled": environment.DeploymentGateEnabled,
			"deployment_gate_check":   environment.DeploymentGateCheck,
			"created_on":              environment.CreatedOn,
			"updated_on":              environment.UpdatedOn,
			"links":                   environment.Links,
		}
	}

	d.Set("environments", environments)
}
