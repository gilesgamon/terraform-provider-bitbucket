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

func dataWorkspaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataWorkspacesRead,
		Schema: map[string]*schema.Schema{
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search query string for workspace names",
			},
			"workspaces": {
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
						"slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_private": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"type": {
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

func dataWorkspacesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG]: params for %s: %v", "dataWorkspacesRead", dumpResourceData(d, dataWorkspaces().Schema))

	url := "2.0/workspaces"

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
		return diag.Errorf("no response returned from workspaces call. Make sure your credentials are accurate.")
	}

	if res.Body == nil {
		return diag.Errorf("error reading workspaces with params (%s): ", dumpResourceData(d, dataWorkspaces().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	workspacesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] workspaces response: %v", workspacesBody)

	var workspacesResponse WorkspacesResponse
	decodeerr := json.Unmarshal(workspacesBody, &workspacesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId("workspaces")
	flattenWorkspaces(&workspacesResponse, d)
	return nil
}

// WorkspacesResponse represents the response from the workspaces API
type WorkspacesResponse struct {
	Values []Workspace `json:"values"`
	Page   int         `json:"page"`
	Size   int         `json:"size"`
	Next   string      `json:"next"`
}

// Workspace represents a workspace in Bitbucket
type Workspace struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Slug      string                 `json:"slug"`
	IsPrivate bool                   `json:"is_private"`
	Type      string                 `json:"type"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the workspaces information
func flattenWorkspaces(c *WorkspacesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	workspaces := make([]interface{}, len(c.Values))
	for i, workspace := range c.Values {
		workspaces[i] = map[string]interface{}{
			"uuid":       workspace.UUID,
			"name":       workspace.Name,
			"slug":       workspace.Slug,
			"is_private": workspace.IsPrivate,
			"type":       workspace.Type,
			"created_on": workspace.CreatedOn,
			"updated_on": workspace.UpdatedOn,
			"links":      workspace.Links,
		}
	}

	d.Set("workspaces", workspaces)
}
