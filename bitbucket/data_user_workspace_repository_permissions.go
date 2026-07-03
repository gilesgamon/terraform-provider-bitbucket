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

func dataUserWorkspaceRepositoryPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataUserWorkspaceRepositoryPermissionsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Query string to narrow down the response.",
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Field by which the results should be sorted.",
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository_full_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataUserWorkspaceRepositoryPermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataUserWorkspaceRepositoryPermissionsRead", dumpResourceData(d, dataUserWorkspaceRepositoryPermissions().Schema))

	url := fmt.Sprintf("2.0/user/workspaces/%s/permissions/repositories", workspace)

	params := make(map[string]string)
	if v, ok := d.GetOk("q"); ok {
		params["q"] = v.(string)
	}
	if v, ok := d.GetOk("sort"); ok {
		params["sort"] = v.(string)
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
	res, err := client.GetAll(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from user workspace repository permissions call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate workspace %s", workspace)
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	body, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)

	var permsResponse UserRepositoryPermissionsResponse
	decodeerr := json.Unmarshal(body, &permsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("user/workspaces/%s/permissions/repositories", workspace))
	d.Set("permissions", flattenUserRepositoryPermissions(permsResponse.Values))
	return nil
}

// UserRepositoryPermissionsResponse represents a paginated list of repository permissions
type UserRepositoryPermissionsResponse struct {
	Values  []UserRepositoryPermission `json:"values"`
	Page    int                        `json:"page"`
	Size    int                        `json:"size"`
	Pagelen int                        `json:"pagelen"`
	Next    string                     `json:"next"`
}

// UserRepositoryPermission represents a user's permission for a given repository
type UserRepositoryPermission struct {
	Permission string `json:"permission"`
	Repository struct {
		UUID     string `json:"uuid"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func flattenUserRepositoryPermissions(items []UserRepositoryPermission) []interface{} {
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"permission":           item.Permission,
			"repository_uuid":      item.Repository.UUID,
			"repository_name":      item.Repository.Name,
			"repository_full_name": item.Repository.FullName,
		}
	}
	return result
}
