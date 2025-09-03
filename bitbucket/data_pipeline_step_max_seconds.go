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

func dataPipelineStepMaxSeconds() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepMaxSecondsRead,
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
			"max_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataPipelineStepMaxSecondsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepMaxSecondsRead", dumpResourceData(d, dataPipelineStepMaxSeconds().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/max-seconds", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step max seconds call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step max seconds with params (%s): ", dumpResourceData(d, dataPipelineStepMaxSeconds().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	maxSecondsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step max seconds response: %v", maxSecondsBody)

	var maxSecondsResponse PipelineStepMaxSecondsResponse
	decodeerr := json.Unmarshal(maxSecondsBody, &maxSecondsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/max-seconds", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepMaxSeconds(&maxSecondsResponse, d)
	return nil
}

// PipelineStepMaxSecondsResponse represents the response from the pipeline step max seconds API
type PipelineStepMaxSecondsResponse struct {
	MaxSeconds int `json:"max_seconds"`
}

// Flattens the pipeline step max seconds information
func flattenPipelineStepMaxSeconds(c *PipelineStepMaxSecondsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("max_seconds", c.MaxSeconds)
}
