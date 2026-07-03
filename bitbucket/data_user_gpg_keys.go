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

func dataUserGpgKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataUserGpgKeysRead,
		Schema: map[string]*schema.Schema{
			"selected_user": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "User UUID or username",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"gpg_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"fingerprint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GPG key fingerprint",
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
				},
			},
		},
	}
}

func dataUserGpgKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	selectedUser := d.Get("selected_user").(string)

	endpoint := fmt.Sprintf("2.0/users/%s/gpg-keys", selectedUser)

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from GPG keys call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate GPG keys for user %s", selectedUser)
	}

	if res.Body == nil {
		return diag.Errorf("error reading GPG keys: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var gpgKeysResponse struct {
		Values []GpgKey `json:"values"`
		Next   string   `json:"next"`
		Size   int      `json:"size"`
		Page   int      `json:"page"`
	}

	if err := json.Unmarshal(body, &gpgKeysResponse); err != nil {
		return diag.FromErr(err)
	}

	var gpgKeys []map[string]interface{}
	for _, key := range gpgKeysResponse.Values {
		keyMap := map[string]interface{}{
			"type":        key.Type,
			"key":         key.Key,
			"fingerprint": key.Fingerprint,
			"created_on":  key.CreatedOn,
		}

		if key.Owner != nil {
			keyMap["owner"] = []map[string]interface{}{
				{
					"display_name": key.Owner.DisplayName,
					"uuid":         key.Owner.UUID,
					"username":     key.Owner.Username,
				},
			}
		}

		if key.Links != nil {
			keyMap["links"] = []map[string]interface{}{
				{
					"self": []map[string]interface{}{
						{
							"href": key.Links.Self.Href,
						},
					},
				},
			}
		}

		gpgKeys = append(gpgKeys, keyMap)
	}

	d.SetId(fmt.Sprintf("gpg-keys-%s", selectedUser))
	d.Set("gpg_keys", gpgKeys)

	log.Printf("[DEBUG] Found %d GPG keys for user %s", len(gpgKeys), selectedUser)

	return nil
}

// GpgKey represents a GPG key
type GpgKey struct {
	Type        string       `json:"type"`
	Key         string       `json:"key"`
	Fingerprint string       `json:"fingerprint"`
	CreatedOn   string       `json:"created_on"`
	Owner       *Account     `json:"owner,omitempty"`
	Links       *GpgKeyLinks `json:"links,omitempty"`
}

// GpgKeyLinks represents GPG key links
type GpgKeyLinks struct {
	Self Link `json:"self"`
}
