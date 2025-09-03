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

func dataPipelineStepAnnotations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepAnnotationsRead,
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
			"annotations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"file_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataPipelineStepAnnotationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepAnnotationsRead", dumpResourceData(d, dataPipelineStepAnnotations().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/annotations", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step annotations call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step annotations with params (%s): ", dumpResourceData(d, dataPipelineStepAnnotations().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	annotationsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step annotations response: %v", annotationsBody)

	var annotationsResponse PipelineStepAnnotationsResponse
	decodeerr := json.Unmarshal(annotationsBody, &annotationsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/annotations", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepAnnotations(&annotationsResponse, d)
	return nil
}

// PipelineStepAnnotationsResponse represents the response from the pipeline step annotations API
type PipelineStepAnnotationsResponse struct {
	Annotations []PipelineStepAnnotation `json:"annotations"`
}

// PipelineStepAnnotation represents an annotation from a pipeline step
type PipelineStepAnnotation struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Line     int    `json:"line"`
	FilePath string `json:"file_path"`
}

// Flattens the pipeline step annotations information
func flattenPipelineStepAnnotations(c *PipelineStepAnnotationsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	annotations := make([]interface{}, len(c.Annotations))
	for i, annotation := range c.Annotations {
		annotations[i] = map[string]interface{}{
			"id":        annotation.ID,
			"type":      annotation.Type,
			"message":   annotation.Message,
			"severity":  annotation.Severity,
			"line":      annotation.Line,
			"file_path": annotation.FilePath,
		}
	}

	d.Set("annotations", annotations)
}
