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

func dataPipelineDeployments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineDeploymentsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deployments": {
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
						"environment": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
					},
				},
			},
		},
	}
}

func dataPipelineDeploymentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineDeploymentsRead", dumpResourceData(d, dataPipelineDeployments().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/deployments", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline deployments call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline deployments for repository %s", repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline deployments with params (%s): ", dumpResourceData(d, dataPipelineDeployments().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	deploymentsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline deployments response: %v", deploymentsBody)

	var deploymentsResponse PipelineDeploymentsResponse
	decodeerr := json.Unmarshal(deploymentsBody, &deploymentsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/deployments", workspace, repoSlug))
	flattenPipelineDeployments(&deploymentsResponse, d)
	return nil
}

// PipelineDeploymentsResponse represents the response from the pipeline deployments API
type PipelineDeploymentsResponse struct {
	Values []PipelineDeployment `json:"values"`
}

// PipelineDeployment represents a pipeline deployment
type PipelineDeployment struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Environment map[string]interface{} `json:"environment"`
	State       string                 `json:"state"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
}

// Flattens the pipeline deployments information
func flattenPipelineDeployments(c *PipelineDeploymentsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	deployments := make([]interface{}, len(c.Values))
	for i, deployment := range c.Values {
		deployments[i] = map[string]interface{}{
			"uuid":        deployment.UUID,
			"name":        deployment.Name,
			"environment": deployment.Environment,
			"state":       deployment.State,
			"created_on":  deployment.CreatedOn,
			"updated_on":  deployment.UpdatedOn,
		}
	}

	d.Set("deployments", deployments)
}
