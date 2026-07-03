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

func dataAddonClientKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAddonClientKeyRead,
		Schema: map[string]*schema.Schema{
			"addon_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The client key of the Connect addon linked to the Forge app installation.",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The raw response body returned by the API.",
			},
		},
	}
}

func dataAddonClientKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	addonKey := d.Get("addon_key").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataAddonClientKeyRead", dumpResourceData(d, dataAddonClientKey().Schema))

	url := fmt.Sprintf("2.0/addon/%s/client-key", addonKey)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from addon client key call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate addon %s client key", addonKey)
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	body, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)

	content := string(body)
	d.Set("content", content)

	// The response may be a bare string or a JSON object containing a client_key field.
	var payload struct {
		ClientKey string `json:"client_key"`
	}
	if err := json.Unmarshal(body, &payload); err == nil && payload.ClientKey != "" {
		d.Set("client_key", payload.ClientKey)
	} else {
		var bare string
		if err := json.Unmarshal(body, &bare); err == nil && bare != "" {
			d.Set("client_key", bare)
		} else {
			d.Set("client_key", content)
		}
	}

	d.SetId(fmt.Sprintf("addon/%s/client-key", addonKey))
	return nil
}
