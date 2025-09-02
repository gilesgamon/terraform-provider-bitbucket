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

	url := fmt.Sprintf("2.0/workspaces/%s/projects", workspace)

	// Build query parameters
	params := make(map[string]string)
	if q, ok := d.GetOk("q"); ok {
		params["q"] = q.(string)
	}

	// Add query parameters to URL
	if len(params) > 0 {
		url += "?"
		first := true
		for key, value := range params {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", key, value)
			first = false
		}
	}

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from projects call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate workspace %s", workspace)
	}

	if res.Body == nil {
		return diag.Errorf("error reading projects with params (%s): ", dumpResourceData(d, dataProjects().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	projectsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] projects response: %v", projectsBody)

	var projectsResponse ProjectsResponse
	decodeerr := json.Unmarshal(projectsBody, &projectsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/projects", workspace))
	flattenProjects(&projectsResponse, d)
	return nil
}

// ProjectsResponse represents the response from the projects API
type ProjectsResponse struct {
	Values []Project `json:"values"`
	Page   int       `json:"page"`
	Size   int       `json:"size"`
	Next   string    `json:"next"`
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
func flattenProjects(c *ProjectsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	projects := make([]interface{}, len(c.Values))
	for i, project := range c.Values {
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
