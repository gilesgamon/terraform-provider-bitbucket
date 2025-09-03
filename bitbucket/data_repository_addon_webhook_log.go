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

func dataRepositoryAddonWebhookLog() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryAddonWebhookLogRead,
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
			"log_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
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
	}
}

func dataRepositoryAddonWebhookLogRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	addonKey := d.Get("addon_key").(string)
	webhookUUID := d.Get("webhook_uuid").(string)
	logUUID := d.Get("log_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryAddonWebhookLogRead", dumpResourceData(d, dataRepositoryAddonWebhookLog().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/addons/%s/webhooks/%s/logs/%s", workspace, repoSlug, addonKey, webhookUUID, logUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository addon webhook log call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate log %s for webhook %s in addon %s for repository %s/%s", logUUID, webhookUUID, addonKey, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository addon webhook log with params (%s): ", dumpResourceData(d, dataRepositoryAddonWebhookLog().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	logBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository addon webhook log response: %v", logBody)

	var logEntry RepositoryAddonWebhookLog
	decodeerr := json.Unmarshal(logBody, &logEntry)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/addons/%s/webhooks/%s/logs/%s", workspace, repoSlug, addonKey, webhookUUID, logUUID))
	flattenRepositoryAddonWebhookLog(&logEntry, d)
	return nil
}

// Flattens the repository addon webhook log information
func flattenRepositoryAddonWebhookLog(c *RepositoryAddonWebhookLog, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("uuid", c.UUID)
	d.Set("level", c.Level)
	d.Set("message", c.Message)
	d.Set("timestamp", c.Timestamp)
	d.Set("details", c.Details)
}
