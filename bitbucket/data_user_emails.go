package bitbucket

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataUserEmails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataUserEmailsRead,
		Schema: map[string]*schema.Schema{
			"emails": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Email address",
						},
						"is_primary": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this is the primary email",
						},
						"is_confirmed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this email is confirmed",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Email type",
						},
					},
				},
			},
		},
	}
}

func dataUserEmailsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	endpoint := "2.0/user/emails"

	res, err := client.GetAll(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from user emails call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate user emails")
	}

	if res.Body == nil {
		return diag.Errorf("error reading user emails: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var emailsResponse struct {
		Values []UserEmailInfo `json:"values"`
		Next   string          `json:"next"`
		Size   int             `json:"size"`
		Page   int             `json:"page"`
	}

	if err := json.Unmarshal(body, &emailsResponse); err != nil {
		return diag.FromErr(err)
	}

	var emails []map[string]interface{}
	for _, email := range emailsResponse.Values {
		emailMap := map[string]interface{}{
			"email":        email.Email,
			"is_primary":   email.IsPrimary,
			"is_confirmed": email.IsConfirmed,
		}
		emails = append(emails, emailMap)
	}

	d.SetId("user-emails")
	d.Set("emails", emails)

	log.Printf("[DEBUG] Found %d emails for current user", len(emails))

	return nil
}

// UserEmailInfo represents a user email
type UserEmailInfo struct {
	Email       string `json:"email"`
	IsPrimary   bool   `json:"is_primary"`
	IsConfirmed bool   `json:"is_confirmed"`
}
