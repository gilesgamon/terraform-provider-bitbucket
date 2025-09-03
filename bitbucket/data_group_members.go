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

func dataGroupMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataGroupMembersRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nickname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_staff": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"account_status": {
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

func dataGroupMembersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	groupSlug := d.Get("group_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataGroupMembersRead", dumpResourceData(d, dataGroupMembers().Schema))

	url := fmt.Sprintf("2.0/workspaces/%s/groups/%s/members", workspace, groupSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from group members call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate group %s in workspace %s", groupSlug, workspace)
	}

	if res.Body == nil {
		return diag.Errorf("error reading group members with params (%s): ", dumpResourceData(d, dataGroupMembers().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	membersBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] group members response: %v", membersBody)

	var membersResponse GroupMembersResponse
	decodeerr := json.Unmarshal(membersBody, &membersResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/groups/%s/members", workspace, groupSlug))
	flattenGroupMembers(&membersResponse, d)
	return nil
}

// GroupMembersResponse represents the response from the group members API
type GroupMembersResponse struct {
	Values []GroupMember `json:"values"`
	Page   int           `json:"page"`
	Size   int           `json:"size"`
	Next   string        `json:"next"`
}

// GroupMember represents a member in a group
type GroupMember struct {
	Username      string                 `json:"username"`
	DisplayName   string                 `json:"display_name"`
	UUID          string                 `json:"uuid"`
	Type          string                 `json:"type"`
	Nickname      string                 `json:"nickname"`
	AccountID     string                 `json:"account_id"`
	CreatedOn     string                 `json:"created_on"`
	IsStaff       bool                   `json:"is_staff"`
	AccountStatus string                 `json:"account_status"`
	Links         map[string]interface{} `json:"links"`
}

// Flattens the group members information
func flattenGroupMembers(c *GroupMembersResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	members := make([]interface{}, len(c.Values))
	for i, member := range c.Values {
		members[i] = map[string]interface{}{
			"username":       member.Username,
			"display_name":   member.DisplayName,
			"uuid":           member.UUID,
			"type":           member.Type,
			"nickname":       member.Nickname,
			"account_id":     member.AccountID,
			"created_on":     member.CreatedOn,
			"is_staff":       member.IsStaff,
			"account_status": member.AccountStatus,
			"links":          member.Links,
		}
	}

	d.Set("members", members)
}
