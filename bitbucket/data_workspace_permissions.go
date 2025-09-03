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

func dataWorkspacePermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataWorkspacePermissionsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
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

func dataWorkspacePermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataWorkspacePermissionsRead", dumpResourceData(d, dataWorkspacePermissions().Schema))

	url := fmt.Sprintf("2.0/workspaces/%s/permissions", workspace)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from workspace permissions call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate workspace %s", workspace)
	}

	if res.Body == nil {
		return diag.Errorf("error reading workspace permissions with params (%s): ", dumpResourceData(d, dataWorkspacePermissions().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	permissionsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] workspace permissions response: %v", permissionsBody)

	var permissionsResponse WorkspacePermissionsResponse
	decodeerr := json.Unmarshal(permissionsBody, &permissionsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/permissions", workspace))
	flattenWorkspacePermissions(&permissionsResponse, d)
	return nil
}

// WorkspacePermissionsResponse represents the response from the workspace permissions API
type WorkspacePermissionsResponse struct {
	Values []WorkspacePermission `json:"values"`
	Page   int                   `json:"page"`
	Size   int                   `json:"size"`
	Next   string                `json:"next"`
}

// WorkspacePermission represents a permission in a workspace
type WorkspacePermission struct {
	Type       string                 `json:"type"`
	User       map[string]interface{} `json:"user"`
	Group      map[string]interface{} `json:"group"`
	Permission string                 `json:"permission"`
	GrantedBy  map[string]interface{} `json:"granted_by"`
	GrantedAt  string                 `json:"granted_at"`
}

// Flattens the workspace permissions information
func flattenWorkspacePermissions(c *WorkspacePermissionsResponse, d *schema.ResourceData) {
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
