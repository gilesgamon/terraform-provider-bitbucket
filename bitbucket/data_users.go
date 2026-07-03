package bitbucket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataUsersRead,
		Schema: map[string]*schema.Schema{
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search query string for usernames or display names",
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nickname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
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

func dataUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG]: params for %s: %v", "dataUsersRead", dumpResourceData(d, dataUsers().Schema))

	params := make(map[string]string)
	if q, ok := d.GetOk("q"); ok {
		params["q"] = q.(string)
	}
	url := "2.0/users" + encodeQueryParams(params)

	client := m.(Clients).httpClient
	rawValues, err := client.GetPaginated(url)
	if err != nil {
		return diag.FromErr(err)
	}

	users := make([]UserListItem, 0, len(rawValues))
	for _, raw := range rawValues {
		var user UserListItem
		if decodeerr := json.Unmarshal(raw, &user); decodeerr != nil {
			return diag.FromErr(decodeerr)
		}
		users = append(users, user)
	}

	d.SetId("users")
	flattenUsers(users, d)
	return nil
}

// UserListItem represents a user in Bitbucket
type UserListItem struct {
	UUID          string                 `json:"uuid"`
	Username      string                 `json:"username"`
	DisplayName   string                 `json:"display_name"`
	Nickname      string                 `json:"nickname"`
	Type          string                 `json:"type"`
	AccountID     string                 `json:"account_id"`
	CreatedOn     string                 `json:"created_on"`
	IsStaff       bool                   `json:"is_staff"`
	AccountStatus string                 `json:"account_status"`
	Links         map[string]interface{} `json:"links"`
}

// Flattens the users information
func flattenUsers(values []UserListItem, d *schema.ResourceData) {
	users := make([]interface{}, len(values))
	for i, user := range values {
		users[i] = map[string]interface{}{
			"uuid":           user.UUID,
			"username":       user.Username,
			"display_name":   user.DisplayName,
			"nickname":       user.Nickname,
			"type":           user.Type,
			"account_id":     user.AccountID,
			"created_on":     user.CreatedOn,
			"is_staff":       user.IsStaff,
			"account_status": user.AccountStatus,
			"links":          user.Links,
		}
	}

	d.Set("users", users)
}
