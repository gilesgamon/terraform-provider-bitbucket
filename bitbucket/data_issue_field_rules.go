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

func dataIssueFieldRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldRulesRead,
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
			"rules": {
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
						"condition": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"action": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
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

func dataIssueFieldRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldRulesRead", dumpResourceData(d, dataIssueFieldRules().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/rules", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field rules call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field rules with params (%s): ", dumpResourceData(d, dataIssueFieldRules().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	rulesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field rules response: %v", rulesBody)

	var rulesResponse IssueFieldRulesResponse
	decodeerr := json.Unmarshal(rulesBody, &rulesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/rules", workspace, repoSlug, fieldUUID))
	flattenIssueFieldRules(&rulesResponse, d)
	return nil
}

// IssueFieldRulesResponse represents the response from the issue field rules API
type IssueFieldRulesResponse struct {
	Values []IssueFieldRule `json:"values"`
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Next   string           `json:"next"`
}

// IssueFieldRule represents a rule for an issue field
type IssueFieldRule struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Condition map[string]interface{} `json:"condition"`
	Action    map[string]interface{} `json:"action"`
	Enabled   bool                   `json:"enabled"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field rules information
func flattenIssueFieldRules(c *IssueFieldRulesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	rules := make([]interface{}, len(c.Values))
	for i, rule := range c.Values {
		rules[i] = map[string]interface{}{
			"uuid":      rule.UUID,
			"name":      rule.Name,
			"type":      rule.Type,
			"condition": rule.Condition,
			"action":    rule.Action,
			"enabled":   rule.Enabled,
			"created_on": rule.CreatedOn,
			"updated_on": rule.UpdatedOn,
			"links":     rule.Links,
		}
	}

	d.Set("rules", rules)
}
