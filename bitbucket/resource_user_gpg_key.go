package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceUserGpgKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGpgKeyCreate,
		ReadContext:   resourceUserGpgKeyRead,
		DeleteContext: resourceUserGpgKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"selected_user": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "User UUID or username",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "GPG public key content",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GPG key type",
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
	}
}

func resourceUserGpgKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	selectedUser := d.Get("selected_user").(string)
	key := d.Get("key").(string)

	gpgKeyRequest := GpgKeyRequest{
		Key: key,
	}

	jsonPayload, err := json.Marshal(gpgKeyRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	endpoint := fmt.Sprintf("2.0/users/%s/gpg-keys", selectedUser)
	res, err := client.Post(endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from GPG key creation")
	}

	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		return diag.Errorf("failed to create GPG key: %s", string(body))
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

	d.SetId(fmt.Sprintf("%s/%s", selectedUser, gpgKey.Fingerprint))

	log.Printf("[DEBUG] Created GPG key: %s for user %s", gpgKey.Fingerprint, selectedUser)

	return resourceUserGpgKeyRead(ctx, d, m)
}

func resourceUserGpgKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	selectedUser, fingerprint, err := userGpgKeyId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	endpoint := fmt.Sprintf("2.0/users/%s/gpg-keys/%s", selectedUser, fingerprint)
	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from GPG key call")
	}

	if res.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] GPG key (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
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

	d.Set("selected_user", selectedUser)
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

func resourceUserGpgKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	selectedUser, fingerprint, err := userGpgKeyId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	endpoint := fmt.Sprintf("2.0/users/%s/gpg-keys/%s", selectedUser, fingerprint)
	res, err := client.Delete(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from GPG key deletion")
	}

	if res.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(res.Body)
		return diag.Errorf("failed to delete GPG key: %s", string(body))
	}

	log.Printf("[DEBUG] Deleted GPG key: %s for user %s", fingerprint, selectedUser)

	return nil
}

// Helper functions
func userGpgKeyId(id string) (selectedUser, fingerprint string, err error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected format of ID (%q), expected USER/FINGERPRINT", id)
	}
	return parts[0], parts[1], nil
}

// GpgKeyRequest represents the request payload for creating GPG keys
type GpgKeyRequest struct {
	Key string `json:"key"`
}
