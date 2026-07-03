package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataGroupsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search query string for group names",
			},
			"groups": {
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
						"slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_private": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
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

func dataGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataGroupsRead", dumpResourceData(d, dataGroups().Schema))

	params := make(map[string]string)
	if q, ok := d.GetOk("q"); ok {
		params["q"] = q.(string)
	}
	url := fmt.Sprintf("2.0/workspaces/%s/groups", workspace) + encodeQueryParams(params)

	client := m.(Clients).httpClient
	rawValues, err := client.GetPaginated(url)
	if err != nil {
		return diag.FromErr(err)
	}

	groups := make([]GroupData, 0, len(rawValues))
	for _, raw := range rawValues {
		var group GroupData
		if decodeerr := json.Unmarshal(raw, &group); decodeerr != nil {
			return diag.FromErr(decodeerr)
		}
		groups = append(groups, group)
	}

	d.SetId(fmt.Sprintf("%s/groups", workspace))
	flattenGroups(groups, d)
	return nil
}

// GroupData represents a group in a workspace from data source
type GroupData struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Description string                 `json:"description"`
	IsPrivate   bool                   `json:"is_private"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Workspace   map[string]interface{} `json:"workspace"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the groups information
func flattenGroups(values []GroupData, d *schema.ResourceData) {
	groups := make([]interface{}, len(values))
	for i, group := range values {
		groups[i] = map[string]interface{}{
			"uuid":        group.UUID,
			"name":        group.Name,
			"slug":        group.Slug,
			"description": group.Description,
			"is_private":  group.IsPrivate,
			"created_on":  group.CreatedOn,
			"updated_on":  group.UpdatedOn,
			"workspace":   group.Workspace,
			"links":       group.Links,
		}
	}

	d.Set("groups", groups)
}
