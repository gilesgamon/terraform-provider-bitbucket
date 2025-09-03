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

func dataIssueFieldConditions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldConditionsRead,
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
			"conditions": {
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
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
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

func dataIssueFieldConditionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldConditionsRead", dumpResourceData(d, dataIssueFieldConditions().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/conditions", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field conditions call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field conditions with params (%s): ", dumpResourceData(d, dataIssueFieldConditions().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	conditionsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field conditions response: %v", conditionsBody)

	var conditionsResponse IssueFieldConditionsResponse
	decodeerr := json.Unmarshal(conditionsBody, &conditionsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/conditions", workspace, repoSlug, fieldUUID))
	flattenIssueFieldConditions(&conditionsResponse, d)
	return nil
}

// IssueFieldConditionsResponse represents the response from the issue field conditions API
type IssueFieldConditionsResponse struct {
	Values []IssueFieldCondition `json:"values"`
	Page   int                   `json:"page"`
	Size   int                   `json:"size"`
	Next   string                `json:"next"`
}

// IssueFieldCondition represents a condition for an issue field
type IssueFieldCondition struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Operator  string                 `json:"operator"`
	Value     string                 `json:"value"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field conditions information
func flattenIssueFieldConditions(c *IssueFieldConditionsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	conditions := make([]interface{}, len(c.Values))
	for i, condition := range c.Values {
		conditions[i] = map[string]interface{}{
			"uuid":       condition.UUID,
			"name":       condition.Name,
			"type":       condition.Type,
			"operator":   condition.Operator,
			"value":      condition.Value,
			"created_on": condition.CreatedOn,
			"updated_on": condition.UpdatedOn,
			"links":      condition.Links,
		}
	}

	d.Set("conditions", conditions)
}
