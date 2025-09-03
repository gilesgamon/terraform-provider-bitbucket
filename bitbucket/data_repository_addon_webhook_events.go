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

func dataRepositoryAddonWebhookEvents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryAddonWebhookEventsRead,
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
			"events": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}

func dataRepositoryAddonWebhookEventsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	addonKey := d.Get("addon_key").(string)
	webhookUUID := d.Get("webhook_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryAddonWebhookEventsRead", dumpResourceData(d, dataRepositoryAddonWebhookEvents().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/addons/%s/webhooks/%s/events", workspace, repoSlug, addonKey, webhookUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository addon webhook events call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate webhook %s for addon %s in repository %s/%s", webhookUUID, addonKey, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository addon webhook events with params (%s): ", dumpResourceData(d, dataRepositoryAddonWebhookEvents().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	eventsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository addon webhook events response: %v", eventsBody)

	var eventsResponse RepositoryAddonWebhookEventsResponse
	decodeerr := json.Unmarshal(eventsBody, &eventsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/addons/%s/webhooks/%s/events", workspace, repoSlug, addonKey, webhookUUID))
	flattenRepositoryAddonWebhookEvents(&eventsResponse, d)
	return nil
}

// RepositoryAddonWebhookEventsResponse represents the response from the repository addon webhook events API
type RepositoryAddonWebhookEventsResponse struct {
	Values []RepositoryAddonWebhookEvent `json:"values"`
	Page   int                           `json:"page"`
	Size   int                           `json:"size"`
	Next   string                        `json:"next"`
}

// RepositoryAddonWebhookEvent represents a webhook event
type RepositoryAddonWebhookEvent struct {
	UUID      string                 `json:"uuid"`
	EventType string                 `json:"event_type"`
	Status    string                 `json:"status"`
	Request   map[string]interface{} `json:"request"`
	Response  map[string]interface{} `json:"response"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the repository addon webhook events information
func flattenRepositoryAddonWebhookEvents(c *RepositoryAddonWebhookEventsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	events := make([]interface{}, len(c.Values))
	for i, event := range c.Values {
		events[i] = map[string]interface{}{
			"uuid":       event.UUID,
			"event_type": event.EventType,
			"status":     event.Status,
			"request":    event.Request,
			"response":   event.Response,
			"created_on": event.CreatedOn,
			"updated_on": event.UpdatedOn,
			"links":      event.Links,
		}
	}

	d.Set("events", events)
}
