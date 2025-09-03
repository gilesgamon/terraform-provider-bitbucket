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

func dataPipelineStepState() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepStateRead,
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"started_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"completed_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataPipelineStepStateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepStateRead", dumpResourceData(d, dataPipelineStepState().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/state", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step state call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step state with params (%s): ", dumpResourceData(d, dataPipelineStepState().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	stateBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step state response: %v", stateBody)

	var stateResponse PipelineStepStateResponse
	decodeerr := json.Unmarshal(stateBody, &stateResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/state", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepState(&stateResponse, d)
	return nil
}

// PipelineStepStateResponse represents the response from the pipeline step state API
type PipelineStepStateResponse struct {
	State        string `json:"state"`
	Name         string `json:"name"`
	StartedOn    string `json:"started_on"`
	CompletedOn  string `json:"completed_on"`
}

// Flattens the pipeline step state information
func flattenPipelineStepState(c *PipelineStepStateResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("state", c.State)
	d.Set("name", c.Name)
	d.Set("started_on", c.StartedOn)
	d.Set("completed_on", c.CompletedOn)
}
