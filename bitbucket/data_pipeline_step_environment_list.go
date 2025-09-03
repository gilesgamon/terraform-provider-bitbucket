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

func dataPipelineStepEnvironmentList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepEnvironmentListRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"step_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataPipelineStepEnvironmentListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepEnvironmentListRead", dumpResourceData(d, dataPipelineStepEnvironmentList().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/environment-list", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step environment list call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step environment list with params (%s): ", dumpResourceData(d, dataPipelineStepEnvironmentList().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	environmentListBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step environment list response: %v", environmentListBody)

	var environmentListResponse PipelineStepEnvironmentListResponse
	decodeerr := json.Unmarshal(environmentListBody, &environmentListResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/environment-list", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepEnvironmentList(&environmentListResponse, d)
	return nil
}

// PipelineStepEnvironmentListResponse represents the response from the pipeline step environment list API
type PipelineStepEnvironmentListResponse struct {
	EnvironmentList []PipelineStepEnvironmentItem `json:"environment_list"`
}

// PipelineStepEnvironmentItem represents an environment item from a pipeline step
type PipelineStepEnvironmentItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// Flattens the pipeline step environment list information
func flattenPipelineStepEnvironmentList(c *PipelineStepEnvironmentListResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	environmentList := make([]interface{}, len(c.EnvironmentList))
	for i, env := range c.EnvironmentList {
		environmentList[i] = map[string]interface{}{
			"key":   env.Key,
			"value": env.Value,
			"type":  env.Type,
		}
	}

	d.Set("environment_list", environmentList)
}
