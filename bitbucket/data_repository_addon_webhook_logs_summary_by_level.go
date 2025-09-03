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

func dataRepositoryAddonWebhookLogsSummaryByLevel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryAddonWebhookLogsSummaryByLevelRead,
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
			"level": {
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

func dataRepositoryAddonWebhookLogsSummaryByLevelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	addonKey := d.Get("addon_key").(string)
	webhookUUID := d.Get("webhook_uuid").(string)
	level := d.Get("level").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryAddonWebhookLogsSummaryByLevelRead", dumpResourceData(d, dataRepositoryAddonWebhookLogsSummaryByLevel().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/addons/%s/webhooks/%s/logs/summary/%s", workspace, repoSlug, addonKey, webhookUUID, level)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository addon webhook logs summary by level call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate webhook %s for addon %s in repository %s/%s", webhookUUID, addonKey, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository addon webhook logs summary by level with params (%s): ", dumpResourceData(d, dataRepositoryAddonWebhookLogsSummaryByLevel().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	summaryBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository addon webhook logs summary by level response: %v", summaryBody)

	var logsSummary RepositoryAddonWebhookLogsSummary
	decodeerr := json.Unmarshal(summaryBody, &logsSummary)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/addons/%s/webhooks/%s/logs/summary/%s", workspace, repoSlug, addonKey, webhookUUID, level))
	flattenRepositoryAddonWebhookLogsSummary(&logsSummary, d)
	return nil
}
