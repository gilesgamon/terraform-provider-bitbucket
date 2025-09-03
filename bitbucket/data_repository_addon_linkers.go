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

func dataRepositoryAddonLinkers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryAddonLinkersRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"linkers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key": {
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
						"application": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func dataRepositoryAddonLinkersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryAddonLinkersRead", dumpResourceData(d, dataRepositoryAddonLinkers().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/addons", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository addon linkers call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository addon linkers with params (%s): ", dumpResourceData(d, dataRepositoryAddonLinkers().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	linkersBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository addon linkers response: %v", linkersBody)

	var linkersResponse RepositoryAddonLinkersResponse
	decodeerr := json.Unmarshal(linkersBody, &linkersResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/addons", workspace, repoSlug))
	flattenRepositoryAddonLinkers(&linkersResponse, d)
	return nil
}

// RepositoryAddonLinkersResponse represents the response from the repository addon linkers API
type RepositoryAddonLinkersResponse struct {
	Values []RepositoryAddonLinker `json:"values"`
	Page   int                     `json:"page"`
	Size   int                     `json:"size"`
	Next   string                  `json:"next"`
}

// RepositoryAddonLinker represents an addon linker
type RepositoryAddonLinker struct {
	UUID        string                 `json:"uuid"`
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Key         string                 `json:"key"`
	Description string                 `json:"description"`
	Vendor      map[string]interface{} `json:"vendor"`
	Application map[string]interface{} `json:"application"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the repository addon linkers information
func flattenRepositoryAddonLinkers(c *RepositoryAddonLinkersResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	linkers := make([]interface{}, len(c.Values))
	for i, linker := range c.Values {
		linkers[i] = map[string]interface{}{
			"uuid":        linker.UUID,
			"id":          linker.ID,
			"name":        linker.Name,
			"key":         linker.Key,
			"description": linker.Description,
			"vendor":      linker.Vendor,
			"application": linker.Application,
			"links":       linker.Links,
		}
	}

	d.Set("linkers", linkers)
}
