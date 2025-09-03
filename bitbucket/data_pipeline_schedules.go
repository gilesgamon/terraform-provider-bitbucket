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

func dataPipelineSchedules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineSchedulesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cron_pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target": {
							Type:     schema.TypeMap,
							Computed: true,
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
				},
			},
		},
	}
}

func dataPipelineSchedulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineSchedulesRead", dumpResourceData(d, dataPipelineSchedules().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/schedules", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline schedules call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline schedules for repository %s", repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline schedules with params (%s): ", dumpResourceData(d, dataPipelineSchedules().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	schedulesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline schedules response: %v", schedulesBody)

	var schedulesResponse PipelineSchedulesResponse
	decodeerr := json.Unmarshal(schedulesBody, &schedulesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines_config/schedules", workspace, repoSlug))
	flattenPipelineSchedules(&schedulesResponse, d)
	return nil
}

// PipelineSchedulesResponse represents the response from the pipeline schedules API
type PipelineSchedulesResponse struct {
	Values []PipelineSchedule `json:"values"`
}

// PipelineSchedule represents a pipeline schedule
type PipelineSchedule struct {
	UUID        string                 `json:"uuid"`
	Enabled     bool                   `json:"enabled"`
	CronPattern string                 `json:"cron_pattern"`
	Target      map[string]interface{} `json:"target"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
}

// Flattens the pipeline schedules information
func flattenPipelineSchedules(c *PipelineSchedulesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	schedules := make([]interface{}, len(c.Values))
	for i, schedule := range c.Values {
		schedules[i] = map[string]interface{}{
			"uuid":         schedule.UUID,
			"enabled":      schedule.Enabled,
			"cron_pattern": schedule.CronPattern,
			"target":       schedule.Target,
			"created_on":   schedule.CreatedOn,
			"updated_on":   schedule.UpdatedOn,
		}
	}

	d.Set("schedules", schedules)
}
