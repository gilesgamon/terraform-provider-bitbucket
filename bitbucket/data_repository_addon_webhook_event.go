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

func dataRepositoryAddonWebhookEvent() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryAddonWebhookEventRead,
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
			"event_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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
	}
}

func dataRepositoryAddonWebhookEventRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	addonKey := d.Get("addon_key").(string)
	webhookUUID := d.Get("webhook_uuid").(string)
	eventUUID := d.Get("event_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryAddonWebhookEventRead", dumpResourceData(d, dataRepositoryAddonWebhookEvent().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/addons/%s/webhooks/%s/events/%s", workspace, repoSlug, addonKey, webhookUUID, eventUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository addon webhook event call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate event %s for webhook %s in addon %s for repository %s/%s", eventUUID, webhookUUID, addonKey, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository addon webhook event with params (%s): ", dumpResourceData(d, dataRepositoryAddonWebhookEvent().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	eventBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository addon webhook event response: %v", eventBody)

	var event RepositoryAddonWebhookEvent
	decodeerr := json.Unmarshal(eventBody, &event)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/addons/%s/webhooks/%s/events/%s", workspace, repoSlug, addonKey, webhookUUID, eventUUID))
	flattenRepositoryAddonWebhookEvent(&event, d)
	return nil
}

// Flattens the repository addon webhook event information
func flattenRepositoryAddonWebhookEvent(c *RepositoryAddonWebhookEvent, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("uuid", c.UUID)
	d.Set("event_type", c.EventType)
	d.Set("status", c.Status)
	d.Set("request", c.Request)
	d.Set("response", c.Response)
	d.Set("created_on", c.CreatedOn)
	d.Set("updated_on", c.UpdatedOn)
	d.Set("links", c.Links)
}
