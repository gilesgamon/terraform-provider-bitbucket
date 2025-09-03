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

func dataPipelineStepArtifactsList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepArtifactsListRead,
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
			"artifacts_list": {
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

func dataPipelineStepArtifactsListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepArtifactsListRead", dumpResourceData(d, dataPipelineStepArtifactsList().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/artifacts-list", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step artifacts list call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step artifacts list with params (%s): ", dumpResourceData(d, dataPipelineStepArtifactsList().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	artifactsListBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step artifacts list response: %v", artifactsListBody)

	var artifactsListResponse PipelineStepArtifactsListResponse
	decodeerr := json.Unmarshal(artifactsListBody, &artifactsListResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/artifacts-list", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepArtifactsList(&artifactsListResponse, d)
	return nil
}

// PipelineStepArtifactsListResponse represents the response from the pipeline step artifacts list API
type PipelineStepArtifactsListResponse struct {
	ArtifactsList []PipelineStepArtifactItem `json:"artifacts_list"`
}

// PipelineStepArtifactItem represents an artifact item from a pipeline step
type PipelineStepArtifactItem struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`
	Path string `json:"path"`
}

// Flattens the pipeline step artifacts list information
func flattenPipelineStepArtifactsList(c *PipelineStepArtifactsListResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	artifactsList := make([]interface{}, len(c.ArtifactsList))
	for i, artifact := range c.ArtifactsList {
		artifactsList[i] = map[string]interface{}{
			"name": artifact.Name,
			"type": artifact.Type,
			"size": artifact.Size,
			"path": artifact.Path,
		}
	}

	d.Set("artifacts_list", artifactsList)
}
