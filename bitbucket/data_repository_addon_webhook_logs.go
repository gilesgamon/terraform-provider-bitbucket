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

func dataRepositoryAddonWebhookLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryAddonWebhookLogsRead,
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
			"webhook_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"details": {
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

func dataRepositoryAddonWebhookLogsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	addonKey := d.Get("addon_key").(string)
	webhookUUID := d.Get("webhook_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryAddonWebhookLogsRead", dumpResourceData(d, dataRepositoryAddonWebhookLogs().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/addons/%s/webhooks/%s/logs", workspace, repoSlug, addonKey, webhookUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository addon webhook logs call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate webhook %s for addon %s in repository %s/%s", webhookUUID, addonKey, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository addon webhook logs with params (%s): ", dumpResourceData(d, dataRepositoryAddonWebhookLogs().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	logsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository addon webhook logs response: %v", logsBody)

	var logsResponse RepositoryAddonWebhookLogsResponse
	decodeerr := json.Unmarshal(logsBody, &logsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/addons/%s/webhooks/%s/logs", workspace, repoSlug, addonKey, webhookUUID))
	flattenRepositoryAddonWebhookLogs(&logsResponse, d)
	return nil
}

// RepositoryAddonWebhookLogsResponse represents the response from the repository addon webhook logs API
type RepositoryAddonWebhookLogsResponse struct {
	Values []RepositoryAddonWebhookLog `json:"values"`
	Page   int                         `json:"page"`
	Size   int                         `json:"size"`
	Next   string                      `json:"next"`
}

// RepositoryAddonWebhookLog represents a webhook log entry
type RepositoryAddonWebhookLog struct {
	UUID      string                 `json:"uuid"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Details   map[string]interface{} `json:"details"`
}

// Flattens the repository addon webhook logs information
func flattenRepositoryAddonWebhookLogs(c *RepositoryAddonWebhookLogsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	logs := make([]interface{}, len(c.Values))
	for i, logEntry := range c.Values {
		logs[i] = map[string]interface{}{
			"uuid":      logEntry.UUID,
			"level":     logEntry.Level,
			"message":   logEntry.Message,
			"timestamp": logEntry.Timestamp,
			"details":   logEntry.Details,
		}
	}

	d.Set("logs", logs)
}
