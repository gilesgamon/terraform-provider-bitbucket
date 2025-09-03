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

func dataIssueFieldAddonWebhookLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldAddonWebhookLogsRead,
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
			"addon_uuid": {
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
						"webhook": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"request": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"response": {
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

func dataIssueFieldAddonWebhookLogsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)
	addonUUID := d.Get("addon_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldAddonWebhookLogsRead", dumpResourceData(d, dataIssueFieldAddonWebhookLogs().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/addons/%s/webhook-logs", workspace, repoSlug, fieldUUID, addonUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field addon webhook logs call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field addon %s in repository %s/%s", addonUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field addon webhook logs with params (%s): ", dumpResourceData(d, dataIssueFieldAddonWebhookLogs().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	logsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field addon webhook logs response: %v", logsBody)

	var logsResponse IssueFieldAddonWebhookLogsResponse
	decodeerr := json.Unmarshal(logsBody, &logsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/addons/%s/webhook-logs", workspace, repoSlug, fieldUUID, addonUUID))
	flattenIssueFieldAddonWebhookLogs(&logsResponse, d)
	return nil
}

// IssueFieldAddonWebhookLogsResponse represents the response from the issue field addon webhook logs API
type IssueFieldAddonWebhookLogsResponse struct {
	Values []IssueFieldAddonWebhookLog `json:"values"`
	Page   int                         `json:"page"`
	Size   int                         `json:"size"`
	Next   string                      `json:"next"`
}

// IssueFieldAddonWebhookLog represents a webhook log for an issue field addon
type IssueFieldAddonWebhookLog struct {
	UUID      string                 `json:"uuid"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Webhook   map[string]interface{} `json:"webhook"`
	Request   map[string]interface{} `json:"request"`
	Response  map[string]interface{} `json:"response"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field addon webhook logs information
func flattenIssueFieldAddonWebhookLogs(c *IssueFieldAddonWebhookLogsResponse, d *schema.ResourceData) {
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
			"webhook":   logEntry.Webhook,
			"request":   logEntry.Request,
			"response":  logEntry.Response,
			"links":     logEntry.Links,
		}
	}

	d.Set("logs", logs)
}
