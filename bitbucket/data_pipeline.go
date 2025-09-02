package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/DrFaust92/bitbucket-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataPipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pipeline number to retrieve",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"build_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"completed_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trigger": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"display_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"target": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ref_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineNumber := d.Get("pipeline_number").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineRead", dumpResourceData(d, dataPipeline().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s",
		workspace,
		repoSlug,
		pipelineNumber,
	)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline %s in repository %s/%s", pipelineNumber, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline information with params (%s): ", dumpResourceData(d, dataPipeline().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	pipelineBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline response: %v", pipelineBody)

	var pipeline Pipeline
	decodeerr := json.Unmarshal(pipelineBody, &pipeline)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", workspace, repoSlug, pipeline.UUID))
	flattenPipeline(&pipeline, d)
	return nil
}

// Pipeline represents a Bitbucket pipeline
type Pipeline struct {
	UUID        string                 `json:"uuid"`
	BuildNumber int                    `json:"build_number"`
	State       PipelineState          `json:"state"`
	Trigger     PipelineTrigger        `json:"trigger"`
	Repository  bitbucket.Repository   `json:"repository"`
	Target      PipelineTarget         `json:"target"`
	CreatedOn   string                 `json:"created_on"`
	CompletedOn string                 `json:"completed_on"`
	Duration    int                    `json:"duration_in_seconds"`
	Links       map[string]interface{} `json:"links"`
}

// PipelineState represents the state of a pipeline
type PipelineState struct {
	Name       string             `json:"name"`
	Stage      string             `json:"stage"`
	Result     string             `json:"result"`
	Completed  bool               `json:"completed"`
	Successful bool               `json:"successful"`
	Failed     bool               `json:"failed"`
	InProgress bool               `json:"in_progress"`
	Stopped    bool               `json:"stopped"`
	StoppedBy  *bitbucket.Account `json:"stopped_by,omitempty"`
}

// PipelineTrigger represents what triggered the pipeline
type PipelineTrigger struct {
	Type string             `json:"type"`
	User *bitbucket.Account `json:"user,omitempty"`
}

// PipelineTarget represents the target of the pipeline
type PipelineTarget struct {
	Type    string         `json:"type"`
	RefName string         `json:"ref_name"`
	Commit  PipelineCommit `json:"commit"`
}

// PipelineCommit represents a commit in a pipeline
type PipelineCommit struct {
	Hash  string                 `json:"hash"`
	Type  string                 `json:"type"`
	Links map[string]interface{} `json:"links"`
}

// Flattens the pipeline information
func flattenPipeline(p *Pipeline, d *schema.ResourceData) {
	if p == nil {
		return
	}

	d.Set("id", fmt.Sprintf("%d", p.BuildNumber))
	d.Set("build_number", fmt.Sprintf("%d", p.BuildNumber))
	d.Set("state", p.State.Name)
	d.Set("created_on", p.CreatedOn)
	d.Set("completed_on", p.CompletedOn)
	d.Set("trigger", flattenPipelineTrigger(&p.Trigger))
	d.Set("target", flattenPipelineTarget(&p.Target))
}

// Flattens the pipeline trigger information
func flattenPipelineTrigger(t *PipelineTrigger) []interface{} {
	if t == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"type": t.Type,
			"user": flattenPipelineAccount(t.User),
		},
	}
}

// Flattens the pipeline target information
func flattenPipelineTarget(t *PipelineTarget) []interface{} {
	if t == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"type":     t.Type,
			"hash":     t.Commit.Hash,
			"ref_name": t.RefName,
		},
	}
}

// Flattens an account (for user information)
func flattenPipelineAccount(a *bitbucket.Account) []interface{} {
	if a == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"username":     a.Username,
			"display_name": a.DisplayName,
			"uuid":         a.Uuid,
		},
	}
}
