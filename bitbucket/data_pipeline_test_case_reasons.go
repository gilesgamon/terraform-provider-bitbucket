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

func dataPipelineTestCaseReasons() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineTestCaseReasonsRead,
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
			"pipeline_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Pipeline UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"step_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Step UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"test_case_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Test case UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"reasons": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reason type",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reason name",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reason description",
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

func dataPipelineTestCaseReasonsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)
	testCaseUUID := d.Get("test_case_uuid").(string)

	endpoint := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/test_reports/test_cases/%s/test_case_reasons", workspace, repoSlug, pipelineUUID, stepUUID, testCaseUUID)

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from pipeline test case reasons call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s, pipeline %s, step %s, or test case %s", workspace, repoSlug, pipelineUUID, stepUUID, testCaseUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline test case reasons: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var reasonsResponse struct {
		Values []TestCaseReason `json:"values"`
		Next   string           `json:"next"`
		Size   int              `json:"size"`
		Page   int              `json:"page"`
	}

	if err := json.Unmarshal(body, &reasonsResponse); err != nil {
		return diag.FromErr(err)
	}

	var reasons []map[string]interface{}
	for _, reason := range reasonsResponse.Values {
		reasonMap := map[string]interface{}{
			"type":        reason.Type,
			"name":        reason.Name,
			"description": reason.Description,
			"created_on":  reason.CreatedOn,
		}
		reasons = append(reasons, reasonMap)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s/%s", workspace, repoSlug, pipelineUUID, stepUUID, testCaseUUID))
	d.Set("reasons", reasons)

	log.Printf("[DEBUG] Found %d reasons for test case %s in pipeline %s step %s in repository %s/%s", len(reasons), testCaseUUID, pipelineUUID, stepUUID, workspace, repoSlug)

	return nil
}

// TestCaseReason represents a test case reason
type TestCaseReason struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedOn   string `json:"created_on"`
}
