package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWorkspaceHooks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataWorkspaceHooksRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hooks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"events": {
							Type:     schema.TypeList,
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

func dataWorkspaceHooksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)

	url := fmt.Sprintf("2.0/workspaces/%s/hooks", workspace)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from workspace hooks call")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate workspace %s", workspace)
	}

	if res.Body == nil {
		return diag.Errorf("error reading workspace hooks")
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	hooksBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}

	var hooksResponse WorkspaceHooksResponse
	decodeerr := json.Unmarshal(hooksBody, &hooksResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/hooks", workspace))
	flattenWorkspaceHooks(&hooksResponse, d)
	return nil
}

type WorkspaceHooksResponse struct {
	Values []WorkspaceHook `json:"values"`
}

type WorkspaceHook struct {
	UUID        string   `json:"uuid"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Active      bool     `json:"active"`
	Events      []string `json:"events"`
}

func flattenWorkspaceHooks(c *WorkspaceHooksResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	hooks := make([]interface{}, len(c.Values))
	for i, hook := range c.Values {
		hooks[i] = map[string]interface{}{
			"uuid":        hook.UUID,
			"url":         hook.URL,
			"description": hook.Description,
			"active":      hook.Active,
			"events":      hook.Events,
		}
	}

	d.Set("hooks", hooks)
}
