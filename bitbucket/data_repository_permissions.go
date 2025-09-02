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

func dataRepositoryPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryPermissionsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
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
						"repository": {
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

func dataRepositoryPermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryPermissionsRead", dumpResourceData(d, dataRepositoryPermissions().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/permissions", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository permissions call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository permissions with params (%s): ", dumpResourceData(d, dataRepositoryPermissions().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	permissionsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository permissions response: %v", permissionsBody)

	var permissionsResponse RepositoryPermissionsResponse
	decodeerr := json.Unmarshal(permissionsBody, &permissionsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/permissions", workspace, repoSlug))
	flattenRepositoryPermissions(&permissionsResponse, d)
	return nil
}

// RepositoryPermissionsResponse represents the response from the repository permissions API
type RepositoryPermissionsResponse struct {
	Values []RepositoryPermission `json:"values"`
	Page   int                    `json:"page"`
	Size   int                    `json:"size"`
	Next   string                 `json:"next"`
}

// RepositoryPermission represents a permission on a repository
type RepositoryPermission struct {
	Type       string                 `json:"type"`
	User       map[string]interface{} `json:"user"`
	Repository map[string]interface{} `json:"repository"`
	Permission string                 `json:"permission"`
	GrantedBy  map[string]interface{} `json:"granted_by"`
	GrantedAt  string                 `json:"granted_at"`
}

// Flattens the repository permissions information
func flattenRepositoryPermissions(c *RepositoryPermissionsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	permissions := make([]interface{}, len(c.Values))
	for i, perm := range c.Values {
		permissions[i] = map[string]interface{}{
			"type":       perm.Type,
			"user":       perm.User,
			"repository": perm.Repository,
			"permission": perm.Permission,
			"granted_by": perm.GrantedBy,
			"granted_at": perm.GrantedAt,
		}
	}

	d.Set("permissions", permissions)
}
