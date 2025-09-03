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

func dataPipelineStepScript() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepScriptRead,
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
			"script": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataPipelineStepScriptRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepScriptRead", dumpResourceData(d, dataPipelineStepScript().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/script", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step script call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step script with params (%s): ", dumpResourceData(d, dataPipelineStepScript().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	scriptBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step script response: %v", scriptBody)

	var scriptResponse PipelineStepScriptResponse
	decodeerr := json.Unmarshal(scriptBody, &scriptResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/script", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepScript(&scriptResponse, d)
	return nil
}

// PipelineStepScriptResponse represents the response from the pipeline step script API
type PipelineStepScriptResponse struct {
	Script map[string]interface{} `json:"script"`
}

// Flattens the pipeline step script information
func flattenPipelineStepScript(c *PipelineStepScriptResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("script", c.Script)
}
