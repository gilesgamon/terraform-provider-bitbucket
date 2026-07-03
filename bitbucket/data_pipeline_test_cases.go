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

func dataPipelineTestCases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineTestCasesRead,
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
			"test_cases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Test case UUID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Test case name",
						},
						"classname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Test case class name",
						},
						"file": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Test case file",
						},
						"result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Test case result",
						},
						"duration": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Test case duration in seconds",
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

func dataPipelineTestCasesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)
	stepUUID := d.Get("step_uuid").(string)

	endpoint := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/steps/%s/test_reports/test_cases", workspace, repoSlug, pipelineUUID, stepUUID)

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from pipeline test cases call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s, pipeline %s, or step %s", workspace, repoSlug, pipelineUUID, stepUUID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline test cases: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var testCasesResponse struct {
		Values []TestCase `json:"values"`
		Next   string     `json:"next"`
		Size   int        `json:"size"`
		Page   int        `json:"page"`
	}

	if err := json.Unmarshal(body, &testCasesResponse); err != nil {
		return diag.FromErr(err)
	}

	var testCases []map[string]interface{}
	for _, testCase := range testCasesResponse.Values {
		testCaseMap := map[string]interface{}{
			"uuid":       testCase.UUID,
			"name":       testCase.Name,
			"classname":  testCase.Classname,
			"file":       testCase.File,
			"result":     testCase.Result,
			"duration":   testCase.Duration,
			"created_on": testCase.CreatedOn,
		}
		testCases = append(testCases, testCaseMap)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", workspace, repoSlug, pipelineUUID, stepUUID))
	d.Set("test_cases", testCases)

	log.Printf("[DEBUG] Found %d test cases for pipeline %s step %s in repository %s/%s", len(testCases), pipelineUUID, stepUUID, workspace, repoSlug)

	return nil
}

// TestCase represents a pipeline test case
type TestCase struct {
	UUID      string  `json:"uuid"`
	Name      string  `json:"name"`
	Classname string  `json:"classname"`
	File      string  `json:"file"`
	Result    string  `json:"result"`
	Duration  float64 `json:"duration"`
	CreatedOn string  `json:"created_on"`
}
