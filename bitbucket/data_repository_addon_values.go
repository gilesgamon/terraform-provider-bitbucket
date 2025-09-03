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

func dataRepositoryAddonValues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryAddonValuesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addon_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataRepositoryAddonValuesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	addonKey := d.Get("addon_key").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryAddonValuesRead", dumpResourceData(d, dataRepositoryAddonValues().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/addons/%s", workspace, repoSlug, addonKey)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository addon values call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate addon %s for repository %s/%s", addonKey, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository addon values with params (%s): ", dumpResourceData(d, dataRepositoryAddonValues().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	valuesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository addon values response: %v", valuesBody)

	var addonValues RepositoryAddonValues
	decodeerr := json.Unmarshal(valuesBody, &addonValues)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/addons/%s", workspace, repoSlug, addonKey))
	flattenRepositoryAddonValues(&addonValues, d)
	return nil
}

// RepositoryAddonValues represents the response from the repository addon values API
type RepositoryAddonValues struct {
	Values map[string]interface{} `json:"values"`
}

// Flattens the repository addon values information
func flattenRepositoryAddonValues(c *RepositoryAddonValues, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("values", c.Values)
}
