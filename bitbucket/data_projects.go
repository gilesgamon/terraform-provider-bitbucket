package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataProjectsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search query string for project names",
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key": {
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
						"owner": {
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

func dataProjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataProjectsRead", dumpResourceData(d, dataProjects().Schema))

	params := make(map[string]string)
	if q, ok := d.GetOk("q"); ok {
		params["q"] = q.(string)
	}
	url := fmt.Sprintf("2.0/workspaces/%s/projects", workspace) + encodeQueryParams(params)

	client := m.(Clients).httpClient
	rawValues, err := client.GetPaginated(url)
	if err != nil {
		return diag.FromErr(err)
	}

	projects := make([]Project, 0, len(rawValues))
	for _, raw := range rawValues {
		var project Project
		if decodeerr := json.Unmarshal(raw, &project); decodeerr != nil {
			return diag.FromErr(decodeerr)
		}
		projects = append(projects, project)
	}

	d.SetId(fmt.Sprintf("%s/projects", workspace))
	flattenProjects(projects, d)
	return nil
}

// Project represents a project in a workspace
type Project struct {
	UUID        string                 `json:"uuid"`
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	IsPrivate   bool                   `json:"is_private"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Owner       map[string]interface{} `json:"owner"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the projects information
func flattenProjects(values []Project, d *schema.ResourceData) {
	projects := make([]interface{}, len(values))
	for i, project := range values {
		projects[i] = map[string]interface{}{
			"uuid":        project.UUID,
			"key":         project.Key,
			"name":        project.Name,
			"description": project.Description,
			"is_private":  project.IsPrivate,
			"created_on":  project.CreatedOn,
			"updated_on":  project.UpdatedOn,
			"owner":       project.Owner,
			"links":       project.Links,
		}
	}

	d.Set("projects", projects)
}
