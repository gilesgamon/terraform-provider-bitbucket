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

func dataPipelineStepImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepImageRead,
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
			"image": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataPipelineStepImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepImageRead", dumpResourceData(d, dataPipelineStepImage().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/image", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step image call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step image with params (%s): ", dumpResourceData(d, dataPipelineStepImage().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	imageBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step image response: %v", imageBody)

	var imageResponse PipelineStepImageResponse
	decodeerr := json.Unmarshal(imageBody, &imageResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/image", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepImage(&imageResponse, d)
	return nil
}

// PipelineStepImageResponse represents the response from the pipeline step image API
type PipelineStepImageResponse struct {
	Image string `json:"image"`
}

// Flattens the pipeline step image information
func flattenPipelineStepImage(c *PipelineStepImageResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("image", c.Image)
}
