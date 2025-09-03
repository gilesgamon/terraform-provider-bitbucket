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

func dataPipelineEnvironments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineEnvironmentsRead,
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
						"lock": {
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
					},
				},
			},
		},
	}
}

func dataPipelineEnvironmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineEnvironmentsRead", dumpResourceData(d, dataPipelineEnvironments().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/environments", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline environments call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline environments for repository %s", repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline environments with params (%s): ", dumpResourceData(d, dataPipelineEnvironments().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	environmentsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline environments response: %v", environmentsBody)

	var environmentsResponse PipelineEnvironmentsResponse
	decodeerr := json.Unmarshal(environmentsBody, &environmentsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/environments", workspace, repoSlug))
	flattenPipelineEnvironments(&environmentsResponse, d)
	return nil
}

// PipelineEnvironmentsResponse represents the response from the pipeline environments API
type PipelineEnvironmentsResponse struct {
	Values []PipelineEnvironment `json:"values"`
}

// PipelineEnvironment represents a pipeline environment
type PipelineEnvironment struct {
	UUID            string                 `json:"uuid"`
	Name            string                 `json:"name"`
	EnvironmentType string                 `json:"environment_type"`
	Rank            int                    `json:"rank"`
	Lock            map[string]interface{} `json:"lock"`
	CreatedOn       string                 `json:"created_on"`
	UpdatedOn       string                 `json:"updated_on"`
}

// Flattens the pipeline environments information
func flattenPipelineEnvironments(c *PipelineEnvironmentsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	environments := make([]interface{}, len(c.Values))
	for i, env := range c.Values {
		environments[i] = map[string]interface{}{
			"uuid":             env.UUID,
			"name":             env.Name,
			"environment_type": env.EnvironmentType,
			"rank":             env.Rank,
			"lock":             env.Lock,
			"created_on":       env.CreatedOn,
			"updated_on":       env.UpdatedOn,
		}
	}

	d.Set("environments", environments)
}
