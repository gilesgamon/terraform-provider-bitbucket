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

func dataIssueFieldActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldActionsRead,
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
			"actions": {
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
						"field": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func dataIssueFieldActionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldActionsRead", dumpResourceData(d, dataIssueFieldActions().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/actions", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field actions call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field actions with params (%s): ", dumpResourceData(d, dataIssueFieldActions().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	actionsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field actions response: %v", actionsBody)

	var actionsResponse IssueFieldActionsResponse
	decodeerr := json.Unmarshal(actionsBody, &actionsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/actions", workspace, repoSlug, fieldUUID))
	flattenIssueFieldActions(&actionsResponse, d)
	return nil
}

// IssueFieldActionsResponse represents the response from the issue field actions API
type IssueFieldActionsResponse struct {
	Values []IssueFieldAction `json:"values"`
	Page   int                `json:"page"`
	Size   int                `json:"size"`
	Next   string             `json:"next"`
}

// IssueFieldAction represents an action for an issue field
type IssueFieldAction struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Field     map[string]interface{} `json:"field"`
	Value     string                 `json:"value"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field actions information
func flattenIssueFieldActions(c *IssueFieldActionsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	actions := make([]interface{}, len(c.Values))
	for i, action := range c.Values {
		actions[i] = map[string]interface{}{
			"uuid":       action.UUID,
			"name":       action.Name,
			"type":       action.Type,
			"field":      action.Field,
			"value":      action.Value,
			"created_on": action.CreatedOn,
			"updated_on": action.UpdatedOn,
			"links":      action.Links,
		}
	}

	d.Set("actions", actions)
}
