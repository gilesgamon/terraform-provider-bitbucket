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

func resourceRepositoryPipelineRunner() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceRepositoryPipelineRunnerCreate,
		ReadWithoutTimeout:   resourceRepositoryPipelineRunnerRead,
		UpdateWithoutTimeout: resourceRepositoryPipelineRunnerUpdate,
		DeleteWithoutTimeout: resourceRepositoryPipelineRunnerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repo_slug": {
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

func resourceRepositoryPipelineRunnerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	payload, err := json.Marshal(expandPipelineRunnerRequest(d))
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines-config/runners", workspace, repoSlug)
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

	d.SetId(fmt.Sprintf("%s/%s/%s", workspace, repoSlug, runner.UUID))
	setPipelineRunnerAttributes(d, &runner)
	// oauth_client secret is only returned on creation, so capture it now.
	d.Set("oauth_client", stringifyMap(runner.OauthClient))

	return resourceRepositoryPipelineRunnerRead(ctx, d, m)
}

func resourceRepositoryPipelineRunnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace, repoSlug, runnerUUID, err := repositoryRunnerId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines-config/runners/%s", workspace, repoSlug, runnerUUID)
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] Repository Pipeline Runner (%s) not found, removing from state", d.Id())
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
	d.Set("repo_slug", repoSlug)
	setPipelineRunnerAttributes(d, &runner)
	return nil
}

func resourceRepositoryPipelineRunnerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace, repoSlug, runnerUUID, err := repositoryRunnerId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload, err := json.Marshal(expandPipelineRunnerRequest(d))
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines-config/runners/%s", workspace, repoSlug, runnerUUID)
	res, err := client.Put(url, bytes.NewBuffer(payload))
	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	return resourceRepositoryPipelineRunnerRead(ctx, d, m)
}

func resourceRepositoryPipelineRunnerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace, repoSlug, runnerUUID, err := repositoryRunnerId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines-config/runners/%s", workspace, repoSlug, runnerUUID)
	res, err := client.Delete(url)
	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func repositoryRunnerId(id string) (string, string, string, error) {
	parts := strings.SplitN(id, "/", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%q), expected WORKSPACE/REPO-SLUG/RUNNER-UUID", id)
	}
	return parts[0], parts[1], parts[2], nil
}
