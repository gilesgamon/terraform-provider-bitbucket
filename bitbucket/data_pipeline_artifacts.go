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

func dataPipelineArtifacts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineArtifactsRead,
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
						"download_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataPipelineArtifactsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineArtifactsRead", dumpResourceData(d, dataPipelineArtifacts().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/artifacts", workspace, repoSlug, pipelineUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline artifacts call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline %s in repository %s", pipelineUUID, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline artifacts with params (%s): ", dumpResourceData(d, dataPipelineArtifacts().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	artifactsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline artifacts response: %v", artifactsBody)

	var artifactsResponse PipelineArtifactsResponse
	decodeerr := json.Unmarshal(artifactsBody, &artifactsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/artifacts", workspace, repoSlug, pipelineUUID))
	flattenPipelineArtifacts(&artifactsResponse, d)
	return nil
}

// PipelineArtifactsResponse represents the response from the pipeline artifacts API
type PipelineArtifactsResponse struct {
	Values []PipelineArtifact `json:"values"`
}

// PipelineArtifact represents a pipeline artifact
type PipelineArtifact struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
	DownloadURL  string `json:"download_url"`
	CreatedOn    string `json:"created_on"`
}

// Flattens the pipeline artifacts information
func flattenPipelineArtifacts(c *PipelineArtifactsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	artifacts := make([]interface{}, len(c.Values))
	for i, artifact := range c.Values {
		artifacts[i] = map[string]interface{}{
			"name":          artifact.Name,
			"type":          artifact.Type,
			"size":          artifact.Size,
			"download_url":  artifact.DownloadURL,
			"created_on":    artifact.CreatedOn,
		}
	}

	d.Set("artifacts", artifacts)
}
