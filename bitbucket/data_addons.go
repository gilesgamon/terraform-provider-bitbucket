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

func dataAddons() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAddonsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addons": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addon_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vendor": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"app_info": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"installed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
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

func dataAddonsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataAddonsRead", dumpResourceData(d, dataAddons().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/addon", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from addons call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading addons with params (%s): ", dumpResourceData(d, dataAddons().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	addonsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] addons response: %v", addonsBody)

	var addonsResponse AddonsResponse
	decodeerr := json.Unmarshal(addonsBody, &addonsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/addons", workspace, repoSlug))
	flattenAddons(&addonsResponse, d)
	return nil
}

// AddonsResponse represents the response from the addons API
type AddonsResponse struct {
	Values []Addon `json:"values"`
	Page   int     `json:"page"`
	Size   int     `json:"size"`
	Next   string  `json:"next"`
}

// Addon represents an addon in a repository
type Addon struct {
	AddonKey   string                 `json:"addon_key"`
	Name       string                 `json:"name"`
	Description string                 `json:"description"`
	Vendor     map[string]interface{} `json:"vendor"`
	AppInfo    map[string]interface{} `json:"app_info"`
	Installed  bool                   `json:"installed"`
	Enabled    bool                   `json:"enabled"`
	Links      map[string]interface{} `json:"links"`
}

// Flattens the addons information
func flattenAddons(c *AddonsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	addons := make([]interface{}, len(c.Values))
	for i, addon := range c.Values {
		addons[i] = map[string]interface{}{
			"addon_key":   addon.AddonKey,
			"name":        addon.Name,
			"description": addon.Description,
			"vendor":      addon.Vendor,
			"app_info":    addon.AppInfo,
			"installed":   addon.Installed,
			"enabled":     addon.Enabled,
			"links":       addon.Links,
		}
	}

	d.Set("addons", addons)
}
