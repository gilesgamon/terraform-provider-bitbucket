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

func dataRepositoryAddonWebhook() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryAddonWebhookRead,
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
	}
}

func dataRepositoryAddonWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	addonKey := d.Get("addon_key").(string)
	webhookUUID := d.Get("webhook_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryAddonWebhookRead", dumpResourceData(d, dataRepositoryAddonWebhook().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/addons/%s/webhooks/%s", workspace, repoSlug, addonKey, webhookUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository addon webhook call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate webhook %s for addon %s in repository %s/%s", webhookUUID, addonKey, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository addon webhook with params (%s): ", dumpResourceData(d, dataRepositoryAddonWebhook().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	webhookBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository addon webhook response: %v", webhookBody)

	var webhook RepositoryAddonWebhook
	decodeerr := json.Unmarshal(webhookBody, &webhook)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/addons/%s/webhooks/%s", workspace, repoSlug, addonKey, webhookUUID))
	flattenRepositoryAddonWebhook(&webhook, d)
	return nil
}

// Flattens the repository addon webhook information
func flattenRepositoryAddonWebhook(c *RepositoryAddonWebhook, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("uuid", c.UUID)
	d.Set("url", c.URL)
	d.Set("description", c.Description)
	d.Set("subject_type", c.SubjectType)
	d.Set("subject", c.Subject)
	d.Set("active", c.Active)
	d.Set("created_on", c.CreatedOn)
	d.Set("updated_on", c.UpdatedOn)
	d.Set("links", c.Links)
}
