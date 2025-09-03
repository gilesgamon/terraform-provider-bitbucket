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

func dataPipelineStepLogsList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepLogsListRead,
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
			"logs_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"step": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataPipelineStepLogsListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepLogsListRead", dumpResourceData(d, dataPipelineStepLogsList().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/logs-list", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step logs list call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step logs list with params (%s): ", dumpResourceData(d, dataPipelineStepLogsList().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	logsListBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step logs list response: %v", logsListBody)

	var logsListResponse PipelineStepLogsListResponse
	decodeerr := json.Unmarshal(logsListBody, &logsListResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/logs-list", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepLogsList(&logsListResponse, d)
	return nil
}

// PipelineStepLogsListResponse represents the response from the pipeline step logs list API
type PipelineStepLogsListResponse struct {
	LogsList []PipelineStepLogItem `json:"logs_list"`
}

// PipelineStepLogItem represents a log item from a pipeline step
type PipelineStepLogItem struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Step      string `json:"step"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}

// Flattens the pipeline step logs list information
func flattenPipelineStepLogsList(c *PipelineStepLogsListResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	logsList := make([]interface{}, len(c.LogsList))
	for i, log := range c.LogsList {
		logsList[i] = map[string]interface{}{
			"level":      log.Level,
			"message":    log.Message,
			"step":       log.Step,
			"created_on": log.CreatedOn,
			"updated_on": log.UpdatedOn,
		}
	}

	d.Set("logs_list", logsList)
}
