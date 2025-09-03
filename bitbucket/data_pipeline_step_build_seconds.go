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

func dataPipelineStepBuildSeconds() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepBuildSecondsRead,
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
			"build_seconds_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataPipelineStepBuildSecondsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepBuildSecondsRead", dumpResourceData(d, dataPipelineStepBuildSeconds().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/build-seconds", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step build seconds call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step build seconds with params (%s): ", dumpResourceData(d, dataPipelineStepBuildSeconds().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	buildSecondsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step build seconds response: %v", buildSecondsBody)

	var buildSecondsResponse PipelineStepBuildSecondsResponse
	decodeerr := json.Unmarshal(buildSecondsBody, &buildSecondsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/build-seconds", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepBuildSeconds(&buildSecondsResponse, d)
	return nil
}

// PipelineStepBuildSecondsResponse represents the response from the pipeline step build seconds API
type PipelineStepBuildSecondsResponse struct {
	BuildSecondsUsed int `json:"build_seconds_used"`
	MaxSeconds       int `json:"max_seconds"`
}

// Flattens the pipeline step build seconds information
func flattenPipelineStepBuildSeconds(c *PipelineStepBuildSecondsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("build_seconds_used", c.BuildSecondsUsed)
	d.Set("max_seconds", c.MaxSeconds)
}
