package bitbucket

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWorkspaceGpgPublicKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataWorkspaceGpgPublicKeyRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The workspace system GPG public key(s).",
			},
		},
	}
}

func dataWorkspaceGpgPublicKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataWorkspaceGpgPublicKeyRead", dumpResourceData(d, dataWorkspaceGpgPublicKey().Schema))

	url := fmt.Sprintf("2.0/workspaces/%s/settings/gpg/public-key", workspace)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from workspace GPG public key call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate GPG public key for workspace %s", workspace)
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	body, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)

	d.SetId(fmt.Sprintf("%s/settings/gpg/public-key", workspace))
	d.Set("public_key", string(body))
	return nil
}
