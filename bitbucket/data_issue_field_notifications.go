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

func dataIssueFieldNotifications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldNotificationsRead,
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
			"notifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"group": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
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

func dataIssueFieldNotificationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldNotificationsRead", dumpResourceData(d, dataIssueFieldNotifications().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/notifications", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field notifications call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field notifications with params (%s): ", dumpResourceData(d, dataIssueFieldNotifications().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	notificationsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field notifications response: %v", notificationsBody)

	var notificationsResponse IssueFieldNotificationsResponse
	decodeerr := json.Unmarshal(notificationsBody, &notificationsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/notifications", workspace, repoSlug, fieldUUID))
	flattenIssueFieldNotifications(&notificationsResponse, d)
	return nil
}

// IssueFieldNotificationsResponse represents the response from the issue field notifications API
type IssueFieldNotificationsResponse struct {
	Values []IssueFieldNotification `json:"values"`
	Page   int                      `json:"page"`
	Size   int                      `json:"size"`
	Next   string                   `json:"next"`
}

// IssueFieldNotification represents a notification for an issue field
type IssueFieldNotification struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	User      map[string]interface{} `json:"user"`
	Group     map[string]interface{} `json:"group"`
	Email     string                 `json:"email"`
	Enabled   bool                   `json:"enabled"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field notifications information
func flattenIssueFieldNotifications(c *IssueFieldNotificationsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	notifications := make([]interface{}, len(c.Values))
	for i, notification := range c.Values {
		notifications[i] = map[string]interface{}{
			"uuid":       notification.UUID,
			"name":       notification.Name,
			"type":       notification.Type,
			"user":       notification.User,
			"group":      notification.Group,
			"email":      notification.Email,
			"enabled":    notification.Enabled,
			"created_on": notification.CreatedOn,
			"updated_on": notification.UpdatedOn,
			"links":      notification.Links,
		}
	}

	d.Set("notifications", notifications)
}
