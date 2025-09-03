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

func dataPipelineStepArtifacts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepArtifactsRead,
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
			"artifacts": {
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
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataPipelineStepArtifactsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepArtifactsRead", dumpResourceData(d, dataPipelineStepArtifacts().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/artifacts", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step artifacts call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step artifacts with params (%s): ", dumpResourceData(d, dataPipelineStepArtifacts().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	artifactsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step artifacts response: %v", artifactsBody)

	var artifactsResponse PipelineStepArtifactsResponse
	decodeerr := json.Unmarshal(artifactsBody, &artifactsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/artifacts", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepArtifacts(&artifactsResponse, d)
	return nil
}

// PipelineStepArtifactsResponse represents the response from the pipeline step artifacts API
type PipelineStepArtifactsResponse struct {
	Values []PipelineStepArtifact `json:"values"`
}

// PipelineStepArtifact represents a pipeline step artifact
type PipelineStepArtifact struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`
	Path string `json:"path"`
}

// Flattens the pipeline step artifacts information
func flattenPipelineStepArtifacts(c *PipelineStepArtifactsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	artifacts := make([]interface{}, len(c.Values))
	for i, artifact := range c.Values {
		artifacts[i] = map[string]interface{}{
			"name": artifact.Name,
			"type": artifact.Type,
			"size": artifact.Size,
			"path": artifact.Path,
		}
	}

	d.Set("artifacts", artifacts)
}
