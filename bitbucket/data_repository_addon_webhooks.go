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

func dataRepositoryAddonWebhooks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryAddonWebhooksRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addon_key": {
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
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject_type": {
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
						"active": {
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

func dataRepositoryAddonWebhooksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	addonKey := d.Get("addon_key").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryAddonWebhooksRead", dumpResourceData(d, dataRepositoryAddonWebhooks().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/addons/%s/webhooks", workspace, repoSlug, addonKey)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository addon webhooks call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate addon %s for repository %s/%s", addonKey, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository addon webhooks with params (%s): ", dumpResourceData(d, dataRepositoryAddonWebhooks().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	webhooksBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository addon webhooks response: %v", webhooksBody)

	var webhooksResponse RepositoryAddonWebhooksResponse
	decodeerr := json.Unmarshal(webhooksBody, &webhooksResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/addons/%s/webhooks", workspace, repoSlug, addonKey))
	flattenRepositoryAddonWebhooks(&webhooksResponse, d)
	return nil
}

// RepositoryAddonWebhooksResponse represents the response from the repository addon webhooks API
type RepositoryAddonWebhooksResponse struct {
	Values []RepositoryAddonWebhook `json:"values"`
	Page   int                      `json:"page"`
	Size   int                      `json:"size"`
	Next   string                   `json:"next"`
}

// RepositoryAddonWebhook represents an addon webhook
type RepositoryAddonWebhook struct {
	UUID        string                 `json:"uuid"`
	URL         string                 `json:"url"`
	Description string                 `json:"description"`
	SubjectType string                 `json:"subject_type"`
	Subject     map[string]interface{} `json:"subject"`
	Active      bool                   `json:"active"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the repository addon webhooks information
func flattenRepositoryAddonWebhooks(c *RepositoryAddonWebhooksResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	webhooks := make([]interface{}, len(c.Values))
	for i, webhook := range c.Values {
		webhooks[i] = map[string]interface{}{
			"uuid":         webhook.UUID,
			"url":          webhook.URL,
			"description":  webhook.Description,
			"subject_type": webhook.SubjectType,
			"subject":      webhook.Subject,
			"active":       webhook.Active,
			"created_on":   webhook.CreatedOn,
			"updated_on":   webhook.UpdatedOn,
			"links":        webhook.Links,
		}
	}

	d.Set("webhooks", webhooks)
}
