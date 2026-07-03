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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataPipelineScheduleExecutions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineScheduleExecutionsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Workspace slug or UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"repo_slug": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Repository slug or UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"schedule_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Schedule UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"executions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution UUID",
						},
						"pipeline": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Pipeline information",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"build_number": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"created_on": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timestamp",
						},
					},
				},
			},
		},
	}
}

func dataPipelineScheduleExecutionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	scheduleUUID := d.Get("schedule_uuid").(string)

	endpoint := fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/schedules/%s/executions", workspace, repoSlug, scheduleUUID)

	res, err := client.GetAll(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from pipeline schedule executions call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s or schedule %s", workspace, repoSlug, scheduleUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline schedule executions: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var executionsResponse struct {
		Values []ScheduleExecution `json:"values"`
		Next   string              `json:"next"`
		Size   int                 `json:"size"`
		Page   int                 `json:"page"`
	}

	if err := json.Unmarshal(body, &executionsResponse); err != nil {
		return diag.FromErr(err)
	}

	var executions []map[string]interface{}
	for _, execution := range executionsResponse.Values {
		executionMap := map[string]interface{}{
			"uuid":       execution.UUID,
			"created_on": execution.CreatedOn,
		}

		if execution.Pipeline != nil {
			pipeline := []map[string]interface{}{
				{
					"uuid":         execution.Pipeline.UUID,
					"state":        execution.Pipeline.State,
					"build_number": execution.Pipeline.BuildNumber,
					"created_on":   execution.Pipeline.CreatedOn,
				},
			}
			executionMap["pipeline"] = pipeline
		}

		executions = append(executions, executionMap)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", workspace, repoSlug, scheduleUUID))
	d.Set("executions", executions)

	log.Printf("[DEBUG] Found %d executions for schedule %s in repository %s/%s", len(executions), scheduleUUID, workspace, repoSlug)

	return nil
}

// ScheduleExecution represents a pipeline schedule execution
type ScheduleExecution struct {
	UUID      string            `json:"uuid"`
	Pipeline  *SchedulePipeline `json:"pipeline,omitempty"`
	CreatedOn string            `json:"created_on"`
}

// SchedulePipeline represents pipeline information in schedule execution
type SchedulePipeline struct {
	UUID        string `json:"uuid"`
	State       string `json:"state"`
	BuildNumber int    `json:"build_number"`
	CreatedOn   string `json:"created_on"`
}
