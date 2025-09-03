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

func dataIssueFieldPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldPermissionsRead,
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
			"permissions": {
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

func dataIssueFieldPermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldPermissionsRead", dumpResourceData(d, dataIssueFieldPermissions().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/permissions", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field permissions call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field permissions with params (%s): ", dumpResourceData(d, dataIssueFieldPermissions().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	permissionsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field permissions response: %v", permissionsBody)

	var permissionsResponse IssueFieldPermissionsResponse
	decodeerr := json.Unmarshal(permissionsBody, &permissionsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/permissions", workspace, repoSlug, fieldUUID))
	flattenIssueFieldPermissions(&permissionsResponse, d)
	return nil
}

// IssueFieldPermissionsResponse represents the response from the issue field permissions API
type IssueFieldPermissionsResponse struct {
	Values []IssueFieldPermission `json:"values"`
	Page   int                    `json:"page"`
	Size   int                    `json:"size"`
	Next   string                 `json:"next"`
}

// IssueFieldPermission represents a permission for an issue field
type IssueFieldPermission struct {
	UUID       string                 `json:"uuid"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	User       map[string]interface{} `json:"user"`
	Group      map[string]interface{} `json:"group"`
	Permission string                 `json:"permission"`
	CreatedOn  string                 `json:"created_on"`
	UpdatedOn  string                 `json:"updated_on"`
	Links      map[string]interface{} `json:"links"`
}

// Flattens the issue field permissions information
func flattenIssueFieldPermissions(c *IssueFieldPermissionsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	permissions := make([]interface{}, len(c.Values))
	for i, permission := range c.Values {
		permissions[i] = map[string]interface{}{
			"uuid":       permission.UUID,
			"name":       permission.Name,
			"type":       permission.Type,
			"user":       permission.User,
			"group":      permission.Group,
			"permission": permission.Permission,
			"created_on": permission.CreatedOn,
			"updated_on": permission.UpdatedOn,
			"links":      permission.Links,
		}
	}

	d.Set("permissions", permissions)
}
