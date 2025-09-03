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

func dataRepositoryAddonWebhookLogsSummaryByTimeRange() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryAddonWebhookLogsSummaryByTimeRangeRead,
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
			"from": {
				Type:     schema.TypeString,
				Required: true,
			},
			"to": {
				Type:     schema.TypeString,
				Required: true,
			},
			"range": {
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

func dataRepositoryAddonWebhookLogsSummaryByTimeRangeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	addonKey := d.Get("addon_key").(string)
	webhookUUID := d.Get("webhook_uuid").(string)
	from := d.Get("from").(string)
	to := d.Get("to").(string)
	timeRange := d.Get("range").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryAddonWebhookLogsSummaryByTimeRangeRead", dumpResourceData(d, dataRepositoryAddonWebhookLogsSummaryByTimeRange().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/addons/%s/webhooks/%s/logs/summary?from=%s&to=%s&range=%s", workspace, repoSlug, addonKey, webhookUUID, from, to, timeRange)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository addon webhook logs summary by time range call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate webhook %s for addon %s in repository %s/%s", webhookUUID, addonKey, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository addon webhook logs summary by time range with params (%s): ", dumpResourceData(d, dataRepositoryAddonWebhookLogsSummaryByTimeRange().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	summaryBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository addon webhook logs summary by time range response: %v", summaryBody)

	var logsSummary RepositoryAddonWebhookLogsSummary
	decodeerr := json.Unmarshal(summaryBody, &logsSummary)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/addons/%s/webhooks/%s/logs/summary/%s/%s/%s", workspace, repoSlug, addonKey, webhookUUID, from, to, timeRange))
	flattenRepositoryAddonWebhookLogsSummary(&logsSummary, d)
	return nil
}
