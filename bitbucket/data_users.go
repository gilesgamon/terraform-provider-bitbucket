package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	url := "2.0/users"

	// Build query parameters
	params := make(map[string]string)
	if q, ok := d.GetOk("q"); ok {
		params["q"] = q.(string)
	}

	// Add query parameters to URL
	if len(params) > 0 {
		url += "?"
		first := true
		for key, value := range params {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", key, value)
			first = false
		}
	}

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from users call. Make sure your credentials are accurate.")
	}

	if res.Body == nil {
		return diag.Errorf("error reading users with params (%s): ", dumpResourceData(d, dataUsers().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	usersBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] users response: %v", usersBody)

	var usersResponse UsersResponse
	decodeerr := json.Unmarshal(usersBody, &usersResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId("users")
	flattenUsers(&usersResponse, d)
	return nil
}

// UsersResponse represents the response from the users API
type UsersResponse struct {
	Values []UserListItem `json:"values"`
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Next   string         `json:"next"`
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
func flattenUsers(c *UsersResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	users := make([]interface{}, len(c.Values))
	for i, user := range c.Values {
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
