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

func dataRepositoryPipelineSchedules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryPipelineSchedulesRead,
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
						"name": {
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
						"next_run": {
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
						"links": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataRepositoryPipelineSchedulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryPipelineSchedulesRead", dumpResourceData(d, dataRepositoryPipelineSchedules().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/schedules", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository pipeline schedules call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository pipeline schedules with params (%s): ", dumpResourceData(d, dataRepositoryPipelineSchedules().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	schedulesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository pipeline schedules response: %v", schedulesBody)

	var schedulesResponse RepositoryPipelineSchedulesResponse
	decodeerr := json.Unmarshal(schedulesBody, &schedulesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines_config/schedules", workspace, repoSlug))
	flattenRepositoryPipelineSchedules(&schedulesResponse, d)
	return nil
}

// RepositoryPipelineSchedulesResponse represents the response from the repository pipeline schedules API
type RepositoryPipelineSchedulesResponse struct {
	Values []RepositoryPipelineSchedule `json:"values"`
	Page   int                          `json:"page"`
	Size   int                          `json:"size"`
	Next   string                       `json:"next"`
}

// RepositoryPipelineSchedule represents a pipeline schedule
type RepositoryPipelineSchedule struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Enabled     bool                   `json:"enabled"`
	CronPattern string                 `json:"cron_pattern"`
	NextRun     string                 `json:"next_run"`
	Target      map[string]interface{} `json:"target"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the repository pipeline schedules information
func flattenRepositoryPipelineSchedules(c *RepositoryPipelineSchedulesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	schedules := make([]interface{}, len(c.Values))
	for i, schedule := range c.Values {
		schedules[i] = map[string]interface{}{
			"uuid":         schedule.UUID,
			"name":         schedule.Name,
			"enabled":      schedule.Enabled,
			"cron_pattern": schedule.CronPattern,
			"next_run":     schedule.NextRun,
			"target":       schedule.Target,
			"created_on":   schedule.CreatedOn,
			"updated_on":   schedule.UpdatedOn,
			"links":        schedule.Links,
		}
	}

	d.Set("schedules", schedules)
}
