package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataUserWorkspaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataUserWorkspacesRead,
		Schema: map[string]*schema.Schema{
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Field by which the results should be sorted.",
			},
			"administrator": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Only return workspaces where the current user is an administrator.",
			},
			"workspaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"administrator": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slug": {
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
					},
				},
			},
		},
	}
}

func dataUserWorkspacesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG]: params for %s: %v", "dataUserWorkspacesRead", dumpResourceData(d, dataUserWorkspaces().Schema))

	url := "2.0/user/workspaces"

	params := make(map[string]string)
	if v, ok := d.GetOk("sort"); ok {
		params["sort"] = v.(string)
	}
	if v, ok := d.GetOkExists("administrator"); ok {
		params["administrator"] = fmt.Sprintf("%t", v.(bool))
	}

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
		return diag.Errorf("no response returned from user workspaces call. Make sure your credentials are accurate.")
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	workspacesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)

	var workspacesResponse UserWorkspacesResponse
	decodeerr := json.Unmarshal(workspacesBody, &workspacesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId("user/workspaces")
	d.Set("workspaces", flattenUserWorkspaces(workspacesResponse.Values))
	return nil
}

// UserWorkspacesResponse represents a paginated list of workspace access objects
type UserWorkspacesResponse struct {
	Values  []WorkspaceAccess `json:"values"`
	Page    int               `json:"page"`
	Size    int               `json:"size"`
	Pagelen int               `json:"pagelen"`
	Next    string            `json:"next"`
}

// WorkspaceAccess represents a user's permission for a workspace
type WorkspaceAccess struct {
	Administrator bool `json:"administrator"`
	Workspace     struct {
		UUID string `json:"uuid"`
		Slug string `json:"slug"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"workspace"`
}

func flattenUserWorkspaces(items []WorkspaceAccess) []interface{} {
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"administrator": item.Administrator,
			"uuid":          item.Workspace.UUID,
			"slug":          item.Workspace.Slug,
			"name":          item.Workspace.Name,
			"type":          item.Workspace.Type,
		}
	}
	return result
}
