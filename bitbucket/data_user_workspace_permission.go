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

func dataUserWorkspacePermission() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataUserWorkspacePermissionRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_nickname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataUserWorkspacePermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataUserWorkspacePermissionRead", dumpResourceData(d, dataUserWorkspacePermission().Schema))

	url := fmt.Sprintf("2.0/user/workspaces/%s/permission", workspace)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from user workspace permission call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate workspace %s membership", workspace)
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	body, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)

	var membership WorkspaceMembership
	decodeerr := json.Unmarshal(body, &membership)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("user/workspaces/%s/permission", workspace))
	d.Set("user_uuid", membership.User.UUID)
	d.Set("user_nickname", membership.User.Nickname)
	d.Set("workspace_uuid", membership.Workspace.UUID)
	d.Set("workspace_slug", membership.Workspace.Slug)
	return nil
}

// WorkspaceMembership represents a Bitbucket workspace membership linking a user to a workspace
type WorkspaceMembership struct {
	User struct {
		UUID     string `json:"uuid"`
		Nickname string `json:"nickname"`
	} `json:"user"`
	Workspace struct {
		UUID string `json:"uuid"`
		Slug string `json:"slug"`
	} `json:"workspace"`
}
