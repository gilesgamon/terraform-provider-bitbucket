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

func dataIssueFieldWebhooks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldWebhooksRead,
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
			"webhooks": {
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
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"events": {
							Type:     schema.TypeList,
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

func dataIssueFieldWebhooksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldWebhooksRead", dumpResourceData(d, dataIssueFieldWebhooks().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/webhooks", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field webhooks call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field webhooks with params (%s): ", dumpResourceData(d, dataIssueFieldWebhooks().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	webhooksBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field webhooks response: %v", webhooksBody)

	var webhooksResponse IssueFieldWebhooksResponse
	decodeerr := json.Unmarshal(webhooksBody, &webhooksResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/webhooks", workspace, repoSlug, fieldUUID))
	flattenIssueFieldWebhooks(&webhooksResponse, d)
	return nil
}

// IssueFieldWebhooksResponse represents the response from the issue field webhooks API
type IssueFieldWebhooksResponse struct {
	Values []IssueFieldWebhook `json:"values"`
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	Next   string              `json:"next"`
}

// IssueFieldWebhook represents a webhook for an issue field
type IssueFieldWebhook struct {
	UUID      string   `json:"uuid"`
	Name      string   `json:"name"`
	URL       string   `json:"url"`
	Events    []string `json:"events"`
	Enabled   bool     `json:"enabled"`
	CreatedOn string   `json:"created_on"`
	UpdatedOn string   `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field webhooks information
func flattenIssueFieldWebhooks(c *IssueFieldWebhooksResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	webhooks := make([]interface{}, len(c.Values))
	for i, webhook := range c.Values {
		webhooks[i] = map[string]interface{}{
			"uuid":       webhook.UUID,
			"name":       webhook.Name,
			"url":        webhook.URL,
			"events":     webhook.Events,
			"enabled":    webhook.Enabled,
			"created_on": webhook.CreatedOn,
			"updated_on": webhook.UpdatedOn,
			"links":      webhook.Links,
		}
	}

	d.Set("webhooks", webhooks)
}
