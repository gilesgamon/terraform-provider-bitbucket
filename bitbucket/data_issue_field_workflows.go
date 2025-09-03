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

func dataIssueFieldWorkflows() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldWorkflowsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"field_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workflows": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"steps": {
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
									"order": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"actions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"enabled": {
							Type:     schema.TypeBool,
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

func dataIssueFieldWorkflowsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldWorkflowsRead", dumpResourceData(d, dataIssueFieldWorkflows().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/workflows", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field workflows call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field workflows with params (%s): ", dumpResourceData(d, dataIssueFieldWorkflows().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	workflowsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field workflows response: %v", workflowsBody)

	var workflowsResponse IssueFieldWorkflowsResponse
	decodeerr := json.Unmarshal(workflowsBody, &workflowsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/workflows", workspace, repoSlug, fieldUUID))
	flattenIssueFieldWorkflows(&workflowsResponse, d)
	return nil
}

// IssueFieldWorkflowsResponse represents the response from the issue field workflows API
type IssueFieldWorkflowsResponse struct {
	Values []IssueFieldWorkflow `json:"values"`
	Page   int                  `json:"page"`
	Size   int                  `json:"size"`
	Next   string               `json:"next"`
}

// IssueFieldWorkflow represents a workflow for an issue field
type IssueFieldWorkflow struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Steps     []WorkflowStep         `json:"steps"`
	Enabled   bool                   `json:"enabled"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// WorkflowStep represents a step in a workflow
type WorkflowStep struct {
	UUID    string   `json:"uuid"`
	Name    string   `json:"name"`
	Order   int      `json:"order"`
	Actions []string `json:"actions"`
}

// Flattens the issue field workflows information
func flattenIssueFieldWorkflows(c *IssueFieldWorkflowsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	workflows := make([]interface{}, len(c.Values))
	for i, workflow := range c.Values {
		steps := make([]interface{}, len(workflow.Steps))
		for j, step := range workflow.Steps {
			steps[j] = map[string]interface{}{
				"uuid":    step.UUID,
				"name":    step.Name,
				"order":   step.Order,
				"actions": step.Actions,
			}
		}

		workflows[i] = map[string]interface{}{
			"uuid":       workflow.UUID,
			"name":       workflow.Name,
			"type":       workflow.Type,
			"steps":      steps,
			"enabled":    workflow.Enabled,
			"created_on": workflow.CreatedOn,
			"updated_on": workflow.UpdatedOn,
			"links":      workflow.Links,
		}
	}

	d.Set("workflows", workflows)
}
