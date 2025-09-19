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

func dataUserGpgKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataUserGpgKeyRead,
		Schema: map[string]*schema.Schema{
			"selected_user": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "User UUID or username",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"fingerprint": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "GPG key fingerprint",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GPG key type",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GPG public key content",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp",
			},
			"owner": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Key owner",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"links": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "GPG key links",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"self": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataUserGpgKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	selectedUser := d.Get("selected_user").(string)
	fingerprint := d.Get("fingerprint").(string)

	endpoint := fmt.Sprintf("2.0/users/%s/gpg-keys/%s", selectedUser, fingerprint)

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from GPG key call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate GPG key %s for user %s", fingerprint, selectedUser)
	}

	if res.Body == nil {
		return diag.Errorf("error reading GPG key: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var gpgKey GpgKey
	if err := json.Unmarshal(body, &gpgKey); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", selectedUser, fingerprint))
	d.Set("type", gpgKey.Type)
	d.Set("key", gpgKey.Key)
	d.Set("fingerprint", gpgKey.Fingerprint)
	d.Set("created_on", gpgKey.CreatedOn)

	if gpgKey.Owner != nil {
		owner := []map[string]interface{}{
			{
				"display_name": gpgKey.Owner.DisplayName,
				"uuid":         gpgKey.Owner.UUID,
				"username":     gpgKey.Owner.Username,
			},
		}
		d.Set("owner", owner)
	}

	if gpgKey.Links != nil {
		links := []map[string]interface{}{
			{
				"self": []map[string]interface{}{
					{
						"href": gpgKey.Links.Self.Href,
					},
				},
			},
		}
		d.Set("links", links)
	}

	log.Printf("[DEBUG] Retrieved GPG key: %s for user %s", fingerprint, selectedUser)

	return nil
}

