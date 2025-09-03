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

func dataPipelineStepScriptList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepScriptListRead,
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
			"script_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataPipelineStepScriptListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepScriptListRead", dumpResourceData(d, dataPipelineStepScriptList().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/script-list", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step script list call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step script list with params (%s): ", dumpResourceData(d, dataPipelineStepScriptList().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	scriptListBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step script list response: %v", scriptListBody)

	var scriptListResponse PipelineStepScriptListResponse
	decodeerr := json.Unmarshal(scriptListBody, &scriptListResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/script-list", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepScriptList(&scriptListResponse, d)
	return nil
}

// PipelineStepScriptListResponse represents the response from the pipeline step script list API
type PipelineStepScriptListResponse struct {
	ScriptList []PipelineStepScriptItem `json:"script_list"`
}

// PipelineStepScriptItem represents a script item from a pipeline step
type PipelineStepScriptItem struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

// Flattens the pipeline step script list information
func flattenPipelineStepScriptList(c *PipelineStepScriptListResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	scriptList := make([]interface{}, len(c.ScriptList))
	for i, script := range c.ScriptList {
		scriptList[i] = map[string]interface{}{
			"name":    script.Name,
			"type":    script.Type,
			"content": script.Content,
		}
	}

	d.Set("script_list", scriptList)
}
