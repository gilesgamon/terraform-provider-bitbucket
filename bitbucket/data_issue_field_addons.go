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

func dataIssueFieldAddons() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldAddonsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"field_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addons": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_on": {
							Type:     schema.TypeString,
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

func dataIssueFieldAddonsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldAddonsRead", dumpResourceData(d, dataIssueFieldAddons().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/addons", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field addons call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field addons with params (%s): ", dumpResourceData(d, dataIssueFieldAddons().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	addonsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field addons response: %v", addonsBody)

	var addonsResponse IssueFieldAddonsResponse
	decodeerr := json.Unmarshal(addonsBody, &addonsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/addons", workspace, repoSlug, fieldUUID))
	flattenIssueFieldAddons(&addonsResponse, d)
	return nil
}

// IssueFieldAddonsResponse represents the response from the issue field addons API
type IssueFieldAddonsResponse struct {
	Values []IssueFieldAddon `json:"values"`
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Next   string            `json:"next"`
}

// IssueFieldAddon represents an addon for an issue field
type IssueFieldAddon struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Config      map[string]interface{} `json:"config"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the issue field addons information
func flattenIssueFieldAddons(c *IssueFieldAddonsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	addons := make([]interface{}, len(c.Values))
	for i, addon := range c.Values {
		addons[i] = map[string]interface{}{
			"uuid":        addon.UUID,
			"name":        addon.Name,
			"type":        addon.Type,
			"version":     addon.Version,
			"description": addon.Description,
			"enabled":     addon.Enabled,
			"config":      addon.Config,
			"created_on":  addon.CreatedOn,
			"updated_on":  addon.UpdatedOn,
			"links":       addon.Links,
		}
	}

	d.Set("addons", addons)
}
