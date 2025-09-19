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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataUserEmail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataUserEmailRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Email address",
				ValidateFunc: validation.StringIsNotEmpty,
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
	}
}

func dataUserEmailRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	email := d.Get("email").(string)

	endpoint := fmt.Sprintf("2.0/user/emails/%s", email)

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from user email call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate email %s", email)
	}

	if res.Body == nil {
		return diag.Errorf("error reading user email: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var userEmail UserEmailInfo
	if err := json.Unmarshal(body, &userEmail); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(email)
	d.Set("email", userEmail.Email)
	d.Set("is_primary", userEmail.IsPrimary)
	d.Set("is_confirmed", userEmail.IsConfirmed)

	log.Printf("[DEBUG] Retrieved user email: %s", userEmail.Email)

	return nil
}
