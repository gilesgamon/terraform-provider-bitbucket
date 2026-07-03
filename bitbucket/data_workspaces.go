package bitbucket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWorkspaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataWorkspacesRead,
		Schema: map[string]*schema.Schema{
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search query string for workspace names",
			},
			"workspaces": {
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
						"is_private": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
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

func dataWorkspacesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG]: params for %s: %v", "dataWorkspacesRead", dumpResourceData(d, dataWorkspaces().Schema))

	params := make(map[string]string)
	if q, ok := d.GetOk("q"); ok {
		params["q"] = q.(string)
	}
	url := "2.0/workspaces" + encodeQueryParams(params)

	client := m.(Clients).httpClient
	rawValues, err := client.GetPaginated(url)
	if err != nil {
		return diag.FromErr(err)
	}

	workspaces := make([]Workspace, 0, len(rawValues))
	for _, raw := range rawValues {
		var workspace Workspace
		if decodeerr := json.Unmarshal(raw, &workspace); decodeerr != nil {
			return diag.FromErr(decodeerr)
		}
		workspaces = append(workspaces, workspace)
	}

	d.SetId("workspaces")
	flattenWorkspaces(workspaces, d)
	return nil
}

// Workspace represents a workspace in Bitbucket
type Workspace struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Slug      string                 `json:"slug"`
	IsPrivate bool                   `json:"is_private"`
	Type      string                 `json:"type"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the workspaces information
func flattenWorkspaces(values []Workspace, d *schema.ResourceData) {
	workspaces := make([]interface{}, len(values))
	for i, workspace := range values {
		workspaces[i] = map[string]interface{}{
			"uuid":       workspace.UUID,
			"name":       workspace.Name,
			"slug":       workspace.Slug,
			"is_private": workspace.IsPrivate,
			"type":       workspace.Type,
			"created_on": workspace.CreatedOn,
			"updated_on": workspace.UpdatedOn,
			"links":      workspace.Links,
		}
	}

	d.Set("workspaces", workspaces)
}
