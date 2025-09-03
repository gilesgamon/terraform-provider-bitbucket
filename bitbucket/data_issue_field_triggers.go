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

func dataIssueFieldTriggers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldTriggersRead,
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
			"triggers": {
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
						"event": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataIssueFieldTriggersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldTriggersRead", dumpResourceData(d, dataIssueFieldTriggers().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/triggers", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field triggers call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field triggers with params (%s): ", dumpResourceData(d, dataIssueFieldTriggers().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	triggersBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field triggers response: %v", triggersBody)

	var triggersResponse IssueFieldTriggersResponse
	decodeerr := json.Unmarshal(triggersBody, &triggersResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/triggers", workspace, repoSlug, fieldUUID))
	flattenIssueFieldTriggers(&triggersResponse, d)
	return nil
}

// IssueFieldTriggersResponse represents the response from the issue field triggers API
type IssueFieldTriggersResponse struct {
	Values []IssueFieldTrigger `json:"values"`
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	Next   string              `json:"next"`
}

// IssueFieldTrigger represents a trigger for an issue field
type IssueFieldTrigger struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Event     string                 `json:"event"`
	Enabled   bool                   `json:"enabled"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field triggers information
func flattenIssueFieldTriggers(c *IssueFieldTriggersResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	triggers := make([]interface{}, len(c.Values))
	for i, trigger := range c.Values {
		triggers[i] = map[string]interface{}{
			"uuid":       trigger.UUID,
			"name":       trigger.Name,
			"type":       trigger.Type,
			"event":      trigger.Event,
			"enabled":    trigger.Enabled,
			"created_on": trigger.CreatedOn,
			"updated_on": trigger.UpdatedOn,
			"links":      trigger.Links,
		}
	}

	d.Set("triggers", triggers)
}
