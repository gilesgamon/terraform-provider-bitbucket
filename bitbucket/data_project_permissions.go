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

func dataProjectPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataProjectPermissionsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"group": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"permission": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"granted_by": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"granted_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataProjectPermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	projectKey := d.Get("project_key").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataProjectPermissionsRead", dumpResourceData(d, dataProjectPermissions().Schema))

	url := fmt.Sprintf("2.0/workspaces/%s/projects/%s/permissions", workspace, projectKey)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from project permissions call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate project %s in workspace %s", projectKey, workspace)
	}

	if res.Body == nil {
		return diag.Errorf("error reading project permissions with params (%s): ", dumpResourceData(d, dataProjectPermissions().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	permissionsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] project permissions response: %v", permissionsBody)

	var permissionsResponse ProjectPermissionsResponse
	decodeerr := json.Unmarshal(permissionsBody, &permissionsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/projects/%s/permissions", workspace, projectKey))
	flattenProjectPermissions(&permissionsResponse, d)
	return nil
}

// ProjectPermissionsResponse represents the response from the project permissions API
type ProjectPermissionsResponse struct {
	Values []ProjectPermission `json:"values"`
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	Next   string              `json:"next"`
}

// ProjectPermission represents a permission in a project
type ProjectPermission struct {
	Type       string                 `json:"type"`
	User       map[string]interface{} `json:"user"`
	Group      map[string]interface{} `json:"group"`
	Permission string                 `json:"permission"`
	GrantedBy  map[string]interface{} `json:"granted_by"`
	GrantedAt  string                 `json:"granted_at"`
}

// Flattens the project permissions information
func flattenProjectPermissions(c *ProjectPermissionsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	permissions := make([]interface{}, len(c.Values))
	for i, permission := range c.Values {
		permissions[i] = map[string]interface{}{
			"type":       permission.Type,
			"user":       permission.User,
			"group":      permission.Group,
			"permission": permission.Permission,
			"granted_by": permission.GrantedBy,
			"granted_at": permission.GrantedAt,
		}
	}

	d.Set("permissions", permissions)
}
