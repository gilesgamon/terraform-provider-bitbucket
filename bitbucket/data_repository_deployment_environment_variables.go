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

func dataRepositoryDeploymentEnvironmentVariables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryDeploymentEnvironmentVariablesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"variables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secured": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"type": {
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

func dataRepositoryDeploymentEnvironmentVariablesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	environmentUUID := d.Get("environment_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryDeploymentEnvironmentVariablesRead", dumpResourceData(d, dataRepositoryDeploymentEnvironmentVariables().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/environments/%s/variables", workspace, repoSlug, environmentUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository deployment environment variables call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate environment %s for repository %s/%s", environmentUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository deployment environment variables with params (%s): ", dumpResourceData(d, dataRepositoryDeploymentEnvironmentVariables().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	variablesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository deployment environment variables response: %v", variablesBody)

	var variablesResponse RepositoryDeploymentEnvironmentVariablesResponse
	decodeerr := json.Unmarshal(variablesBody, &variablesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/environments/%s/variables", workspace, repoSlug, environmentUUID))
	flattenRepositoryDeploymentEnvironmentVariables(&variablesResponse, d)
	return nil
}

// RepositoryDeploymentEnvironmentVariablesResponse represents the response from the repository deployment environment variables API
type RepositoryDeploymentEnvironmentVariablesResponse struct {
	Values []RepositoryDeploymentEnvironmentVariable `json:"values"`
	Page   int                                        `json:"page"`
	Size   int                                        `json:"size"`
	Next   string                                     `json:"next"`
}

// RepositoryDeploymentEnvironmentVariable represents a deployment environment variable
type RepositoryDeploymentEnvironmentVariable struct {
	UUID      string                 `json:"uuid"`
	Key       string                 `json:"key"`
	Value     string                 `json:"value"`
	Secured   bool                   `json:"secured"`
	Type      string                 `json:"type"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the repository deployment environment variables information
func flattenRepositoryDeploymentEnvironmentVariables(c *RepositoryDeploymentEnvironmentVariablesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	variables := make([]interface{}, len(c.Values))
	for i, variable := range c.Values {
		variables[i] = map[string]interface{}{
			"uuid":       variable.UUID,
			"key":        variable.Key,
			"value":      variable.Value,
			"secured":    variable.Secured,
			"type":       variable.Type,
			"created_on": variable.CreatedOn,
			"updated_on": variable.UpdatedOn,
			"links":      variable.Links,
		}
	}

	d.Set("variables", variables)
}
