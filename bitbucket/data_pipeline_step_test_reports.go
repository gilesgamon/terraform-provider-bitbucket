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

func dataPipelineStepTestReports() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineStepTestReportsRead,
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
			"test_reports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_tests": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"passed_tests": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"failed_tests": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"skipped_tests": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataPipelineStepTestReportsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineStepTestReportsRead", dumpResourceData(d, dataPipelineStepTestReports().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/test-reports", workspace, repoSlug, pipelineUUID, stepUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline step test reports call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline step %s in pipeline %s", stepUUID, pipelineUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline step test reports with params (%s): ", dumpResourceData(d, dataPipelineStepTestReports().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	testReportsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline step test reports response: %v", testReportsBody)

	var testReportsResponse PipelineStepTestReportsResponse
	decodeerr := json.Unmarshal(testReportsBody, &testReportsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/steps/%s/test-reports", workspace, repoSlug, pipelineUUID, stepUUID))
	flattenPipelineStepTestReports(&testReportsResponse, d)
	return nil
}

// PipelineStepTestReportsResponse represents the response from the pipeline step test reports API
type PipelineStepTestReportsResponse struct {
	Values []PipelineStepTestReport `json:"values"`
}

// PipelineStepTestReport represents a pipeline step test report
type PipelineStepTestReport struct {
	Name         string  `json:"name"`
	TotalTests   int     `json:"total_tests"`
	PassedTests  int     `json:"passed_tests"`
	FailedTests  int     `json:"failed_tests"`
	SkippedTests int     `json:"skipped_tests"`
	Duration     float64 `json:"duration"`
}

// Flattens the pipeline step test reports information
func flattenPipelineStepTestReports(c *PipelineStepTestReportsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	testReports := make([]interface{}, len(c.Values))
	for i, report := range c.Values {
		testReports[i] = map[string]interface{}{
			"name":          report.Name,
			"total_tests":   report.TotalTests,
			"passed_tests":  report.PassedTests,
			"failed_tests":  report.FailedTests,
			"skipped_tests": report.SkippedTests,
			"duration":      report.Duration,
		}
	}

	d.Set("test_reports", testReports)
}
