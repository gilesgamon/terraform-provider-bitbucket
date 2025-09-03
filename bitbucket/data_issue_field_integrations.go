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

func dataIssueFieldIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldIntegrationsRead,
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
			"integrations": {
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
						"provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config": {
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

func dataIssueFieldIntegrationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldIntegrationsRead", dumpResourceData(d, dataIssueFieldIntegrations().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/integrations", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field integrations call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field integrations with params (%s): ", dumpResourceData(d, dataIssueFieldIntegrations().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	integrationsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field integrations response: %v", integrationsBody)

	var integrationsResponse IssueFieldIntegrationsResponse
	decodeerr := json.Unmarshal(integrationsBody, &integrationsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/integrations", workspace, repoSlug, fieldUUID))
	flattenIssueFieldIntegrations(&integrationsResponse, d)
	return nil
}

// IssueFieldIntegrationsResponse represents the response from the issue field integrations API
type IssueFieldIntegrationsResponse struct {
	Values []IssueFieldIntegration `json:"values"`
	Page   int                     `json:"page"`
	Size   int                     `json:"size"`
	Next   string                  `json:"next"`
}

// IssueFieldIntegration represents an integration for an issue field
type IssueFieldIntegration struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Provider  string                 `json:"provider"`
	Config    map[string]interface{} `json:"config"`
	Enabled   bool                   `json:"enabled"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field integrations information
func flattenIssueFieldIntegrations(c *IssueFieldIntegrationsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	integrations := make([]interface{}, len(c.Values))
	for i, integration := range c.Values {
		integrations[i] = map[string]interface{}{
			"uuid":       integration.UUID,
			"name":       integration.Name,
			"type":       integration.Type,
			"provider":   integration.Provider,
			"config":     integration.Config,
			"enabled":    integration.Enabled,
			"created_on": integration.CreatedOn,
			"updated_on": integration.UpdatedOn,
			"links":      integration.Links,
		}
	}

	d.Set("integrations", integrations)
}
