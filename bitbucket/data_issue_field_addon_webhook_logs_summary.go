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

func dataIssueFieldAddonWebhookLogsSummary() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldAddonWebhookLogsSummaryRead,
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
			"summary": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataIssueFieldAddonWebhookLogsSummaryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)
	addonUUID := d.Get("addon_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldAddonWebhookLogsSummaryRead", dumpResourceData(d, dataIssueFieldAddonWebhookLogsSummary().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/addons/%s/webhook-logs/summary", workspace, repoSlug, fieldUUID, addonUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field addon webhook logs summary call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field addon %s in repository %s/%s", addonUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field addon webhook logs summary with params (%s): ", dumpResourceData(d, dataIssueFieldAddonWebhookLogsSummary().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	summaryBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field addon webhook logs summary response: %v", summaryBody)

	var summaryResponse IssueFieldAddonWebhookLogsSummaryResponse
	decodeerr := json.Unmarshal(summaryBody, &summaryResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/addons/%s/webhook-logs/summary", workspace, repoSlug, fieldUUID, addonUUID))
	flattenIssueFieldAddonWebhookLogsSummary(&summaryResponse, d)
	return nil
}

// IssueFieldAddonWebhookLogsSummaryResponse represents the response from the issue field addon webhook logs summary API
type IssueFieldAddonWebhookLogsSummaryResponse struct {
	TotalLogs    int                    `json:"total_logs"`
	ErrorLogs    int                    `json:"error_logs"`
	WarningLogs  int                    `json:"warning_logs"`
	InfoLogs     int                    `json:"info_logs"`
	DebugLogs    int                    `json:"debug_logs"`
	LastLogTime  string                 `json:"last_log_time"`
	FirstLogTime string                 `json:"first_log_time"`
	Webhooks     map[string]interface{} `json:"webhooks"`
	Links        map[string]interface{} `json:"links"`
}

// Flattens the issue field addon webhook logs summary information
func flattenIssueFieldAddonWebhookLogsSummary(c *IssueFieldAddonWebhookLogsSummaryResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	summary := map[string]interface{}{
		"total_logs":    c.TotalLogs,
		"error_logs":    c.ErrorLogs,
		"warning_logs":  c.WarningLogs,
		"info_logs":     c.InfoLogs,
		"debug_logs":    c.DebugLogs,
		"last_log_time": c.LastLogTime,
		"first_log_time": c.FirstLogTime,
		"webhooks":      c.Webhooks,
		"links":         c.Links,
	}

	d.Set("summary", summary)
}
