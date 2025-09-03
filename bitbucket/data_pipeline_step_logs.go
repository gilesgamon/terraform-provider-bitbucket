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

func dataPipelineStepLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepLogsRead,
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
			"logs": {
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

func dataPipelineStepLogsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepLogsRead", dumpResourceData(d, dataPipelineStepLogs().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/logs", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step logs call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step logs with params (%s): ", dumpResourceData(d, dataPipelineStepLogs().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	logsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step logs response: %v", logsBody)

	var logsResponse PipelineStepLogsResponse
	decodeerr := json.Unmarshal(logsBody, &logsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/logs", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepLogs(&logsResponse, d)
	return nil
}

// PipelineStepLogsResponse represents the response from the pipeline step logs API
type PipelineStepLogsResponse struct {
	Values []PipelineStepLog `json:"values"`
}

// PipelineStepLog represents a pipeline step log entry
type PipelineStepLog struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Step      string `json:"step"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}

// Flattens the pipeline step logs information
func flattenPipelineStepLogs(c *PipelineStepLogsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	logs := make([]interface{}, len(c.Values))
	for i, logEntry := range c.Values {
		logs[i] = map[string]interface{}{
			"level":      logEntry.Level,
			"message":    logEntry.Message,
			"step":       logEntry.Step,
			"created_on": logEntry.CreatedOn,
			"updated_on": logEntry.UpdatedOn,
		}
	}

	d.Set("logs", logs)
}
