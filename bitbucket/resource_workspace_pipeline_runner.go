package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// pipelineRunnerRequest is the payload used to create or update a runner.
type pipelineRunnerRequest struct {
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
}

func resourceWorkspacePipelineRunner() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceWorkspacePipelineRunnerCreate,
		ReadWithoutTimeout:   resourceWorkspacePipelineRunnerRead,
		UpdateWithoutTimeout: resourceWorkspacePipelineRunnerUpdate,
		DeleteWithoutTimeout: resourceWorkspacePipelineRunnerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"labels": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"oauth_client": {
				Type:      schema.TypeMap,
				Computed:  true,
				Sensitive: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func expandPipelineRunnerRequest(d *schema.ResourceData) *pipelineRunnerRequest {
	labelsSet := d.Get("labels").(*schema.Set).List()
	labels := make([]string, len(labelsSet))
	for i, l := range labelsSet {
		labels[i] = l.(string)
	}
	return &pipelineRunnerRequest{
		Name:   d.Get("name").(string),
		Labels: labels,
	}
}

func resourceWorkspacePipelineRunnerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient
	workspace := d.Get("workspace").(string)

	payload, err := json.Marshal(expandPipelineRunnerRequest(d))
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("2.0/workspaces/%s/pipelines-config/runners", workspace)
	res, err := client.Post(url, bytes.NewBuffer(payload))
	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	body, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}

	var runner PipelineRunner
	if decodeerr := json.Unmarshal(body, &runner); decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s", workspace, runner.UUID))
	setPipelineRunnerAttributes(d, &runner)
	// oauth_client secret is only returned on creation, so capture it now.
	d.Set("oauth_client", stringifyMap(runner.OauthClient))

	return resourceWorkspacePipelineRunnerRead(ctx, d, m)
}

func resourceWorkspacePipelineRunnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace, runnerUUID, err := workspaceRunnerId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("2.0/workspaces/%s/pipelines-config/runners/%s", workspace, runnerUUID)
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] Workspace Pipeline Runner (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	body, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}

	var runner PipelineRunner
	if decodeerr := json.Unmarshal(body, &runner); decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.Set("workspace", workspace)
	setPipelineRunnerAttributes(d, &runner)
	return nil
}

func resourceWorkspacePipelineRunnerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace, runnerUUID, err := workspaceRunnerId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload, err := json.Marshal(expandPipelineRunnerRequest(d))
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("2.0/workspaces/%s/pipelines-config/runners/%s", workspace, runnerUUID)
	res, err := client.Put(url, bytes.NewBuffer(payload))
	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	return resourceWorkspacePipelineRunnerRead(ctx, d, m)
}

func resourceWorkspacePipelineRunnerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace, runnerUUID, err := workspaceRunnerId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("2.0/workspaces/%s/pipelines-config/runners/%s", workspace, runnerUUID)
	res, err := client.Delete(url)
	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// setPipelineRunnerAttributes populates the common runner attributes on the resource data.
func setPipelineRunnerAttributes(d *schema.ResourceData, runner *PipelineRunner) {
	d.Set("name", runner.Name)
	d.Set("uuid", runner.UUID)
	d.Set("labels", runner.Labels)
	d.Set("created_on", runner.CreatedOn)
	d.Set("updated_on", runner.UpdatedOn)

	if runner.State != nil {
		state := map[string]interface{}{
			"status":     runner.State.Status,
			"cordoned":   fmt.Sprintf("%t", runner.State.Cordoned),
			"updated_on": runner.State.UpdatedOn,
		}
		d.Set("state", state)
	}
}

func workspaceRunnerId(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%q), expected WORKSPACE/RUNNER-UUID", id)
	}
	return parts[0], parts[1], nil
}
