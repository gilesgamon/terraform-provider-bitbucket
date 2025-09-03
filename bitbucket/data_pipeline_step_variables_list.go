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

func dataPipelineStepVariablesList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepVariablesListRead,
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
			"variables_list": {
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
						"secured": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataPipelineStepVariablesListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepVariablesListRead", dumpResourceData(d, dataPipelineStepVariablesList().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/variables-list", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step variables list call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step variables list with params (%s): ", dumpResourceData(d, dataPipelineStepVariablesList().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	variablesListBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step variables list response: %v", variablesListBody)

	var variablesListResponse PipelineStepVariablesListResponse
	decodeerr := json.Unmarshal(variablesListBody, &variablesListResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/variables-list", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepVariablesList(&variablesListResponse, d)
	return nil
}

// PipelineStepVariablesListResponse represents the response from the pipeline step variables list API
type PipelineStepVariablesListResponse struct {
	VariablesList []PipelineStepVariableItem `json:"variables_list"`
}

// PipelineStepVariableItem represents a variable item from a pipeline step
type PipelineStepVariableItem struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Type     string `json:"type"`
	Secured  bool   `json:"secured"`
}

// Flattens the pipeline step variables list information
func flattenPipelineStepVariablesList(c *PipelineStepVariablesListResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	variablesList := make([]interface{}, len(c.VariablesList))
	for i, variable := range c.VariablesList {
		variablesList[i] = map[string]interface{}{
			"key":     variable.Key,
			"value":   variable.Value,
			"type":    variable.Type,
			"secured": variable.Secured,
		}
	}

	d.Set("variables_list", variablesList)
}
