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

func dataPipelineSteps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepsRead,
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
			"steps": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
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
						"duration_in_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"script": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func dataPipelineStepsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepsRead", dumpResourceData(d, dataPipelineSteps().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps", workspace, repoSlug, pipelineUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline steps call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline %s in repository %s/%s", pipelineUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline steps with params (%s): ", dumpResourceData(d, dataPipelineSteps().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	stepsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline steps response: %v", stepsBody)

	var stepsResponse PipelineStepsResponse
	decodeerr := json.Unmarshal(stepsBody, &stepsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps", workspace, repoSlug, pipelineUUID))
	flattenPipelineSteps(&stepsResponse, d)
	return nil
}

// PipelineStepsResponse represents the response from the pipeline steps API
type PipelineStepsResponse struct {
	Values []PipelineStep `json:"values"`
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Next   string         `json:"next"`
}

// PipelineStep represents a pipeline step
type PipelineStep struct {
	UUID                string                 `json:"uuid"`
	Name                string                 `json:"name"`
	Type                string                 `json:"type"`
	State               string                 `json:"state"`
	StartedOn           string                 `json:"started_on"`
	CompletedOn         string                 `json:"completed_on"`
	DurationInSeconds   int                    `json:"duration_in_seconds"`
	MaxTime             int                    `json:"max_time"`
	Script              map[string]interface{} `json:"script"`
	Links               map[string]interface{} `json:"links"`
}

// Flattens the pipeline steps information
func flattenPipelineSteps(c *PipelineStepsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	steps := make([]interface{}, len(c.Values))
	for i, step := range c.Values {
		steps[i] = map[string]interface{}{
			"uuid":                  step.UUID,
			"name":                  step.Name,
			"type":                  step.Type,
			"state":                 step.State,
			"started_on":            step.StartedOn,
			"completed_on":          step.CompletedOn,
			"duration_in_seconds":   step.DurationInSeconds,
			"max_time":              step.MaxTime,
			"script":                step.Script,
			"links":                 step.Links,
		}
	}

	d.Set("steps", steps)
}
