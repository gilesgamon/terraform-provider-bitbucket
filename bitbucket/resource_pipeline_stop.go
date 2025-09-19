package bitbucket

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePipelineStop() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineStopCreate,
		ReadContext:   resourcePipelineStopRead,
		DeleteContext: resourcePipelineStopDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Workspace slug or UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"repo_slug": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Repository slug or UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"pipeline_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Pipeline UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"stopped": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the pipeline was successfully stopped",
			},
		},
	}
}

func resourcePipelineStopCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)

	endpoint := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/stopPipeline", workspace, repoSlug, pipelineUUID)
	res, err := client.Post(endpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from pipeline stop")
	}

	if res.StatusCode != http.StatusNoContent {
		return diag.Errorf("failed to stop pipeline: status %d", res.StatusCode)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", workspace, repoSlug, pipelineUUID))
	d.Set("stopped", true)

	log.Printf("[DEBUG] Stopped pipeline: %s in repository %s/%s", pipelineUUID, workspace, repoSlug)

	return nil
}

func resourcePipelineStopRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Pipeline stop is a one-time action, so we don't need to read state
	// Just return the current state
	return nil
}

func resourcePipelineStopDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Pipeline stop is a one-time action, so deletion is a no-op
	// Just remove from state
	d.SetId("")
	return nil
}

