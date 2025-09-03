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

func dataWebhooks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataWebhooksRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"events": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"skip_cert_verification": {
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
						"subject": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func dataWebhooksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataWebhooksRead", dumpResourceData(d, dataWebhooks().Schema))

	url := fmt.Sprintf("2.0/workspaces/%s/hooks", workspace)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from webhooks call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate workspace %s", workspace)
	}

	if res.Body == nil {
		return diag.Errorf("error reading webhooks with params (%s): ", dumpResourceData(d, dataWebhooks().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	webhooksBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] webhooks response: %v", webhooksBody)

	var webhooksResponse WebhooksResponse
	decodeerr := json.Unmarshal(webhooksBody, &webhooksResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/webhooks", workspace))
	flattenWebhooks(&webhooksResponse, d)
	return nil
}

// WebhooksResponse represents the response from the webhooks API
type WebhooksResponse struct {
	Values []Webhook `json:"values"`
	Page   int       `json:"page"`
	Size   int       `json:"size"`
	Next   string    `json:"next"`
}

// Webhook represents a webhook
type Webhook struct {
	UUID                 string                   `json:"uuid"`
	Description          string                   `json:"description"`
	URL                  string                   `json:"url"`
	Active               bool                     `json:"active"`
	Events               []string                 `json:"events"`
	SkipCertVerification bool                     `json:"skip_cert_verification"`
	CreatedOn            string                   `json:"created_on"`
	UpdatedOn            string                   `json:"updated_on"`
	Subject              map[string]interface{}  `json:"subject"`
	Links                map[string]interface{}  `json:"links"`
}

// Flattens the webhooks information
func flattenWebhooks(c *WebhooksResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	webhooks := make([]interface{}, len(c.Values))
	for i, webhook := range c.Values {
		webhooks[i] = map[string]interface{}{
			"uuid":                   webhook.UUID,
			"description":            webhook.Description,
			"url":                    webhook.URL,
			"active":                 webhook.Active,
			"events":                 webhook.Events,
			"skip_cert_verification": webhook.SkipCertVerification,
			"created_on":             webhook.CreatedOn,
			"updated_on":             webhook.UpdatedOn,
			"subject":                webhook.Subject,
			"links":                  webhook.Links,
		}
	}

	d.Set("webhooks", webhooks)
}
