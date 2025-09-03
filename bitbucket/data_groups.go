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

func dataGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataGroupsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search query string for group names",
			},
			"groups": {
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
						"workspace": {
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

func dataGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataGroupsRead", dumpResourceData(d, dataGroups().Schema))

	url := fmt.Sprintf("2.0/workspaces/%s/groups", workspace)

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
		return diag.Errorf("no response returned from groups call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate workspace %s", workspace)
	}

	if res.Body == nil {
		return diag.Errorf("error reading groups with params (%s): ", dumpResourceData(d, dataGroups().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	groupsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] groups response: %v", groupsBody)

	var groupsResponse GroupsResponse
	decodeerr := json.Unmarshal(groupsBody, &groupsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/groups", workspace))
	flattenGroups(&groupsResponse, d)
	return nil
}

// GroupsResponse represents the response from the groups API
type GroupsResponse struct {
	Values []GroupData `json:"values"`
	Page   int         `json:"page"`
	Size   int         `json:"size"`
	Next   string      `json:"next"`
}

// GroupData represents a group in a workspace from data source
type GroupData struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Description string                 `json:"description"`
	IsPrivate   bool                   `json:"is_private"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Workspace   map[string]interface{} `json:"workspace"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the groups information
func flattenGroups(c *GroupsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	groups := make([]interface{}, len(c.Values))
	for i, group := range c.Values {
		groups[i] = map[string]interface{}{
			"uuid":        group.UUID,
			"name":        group.Name,
			"slug":        group.Slug,
			"description": group.Description,
			"is_private":  group.IsPrivate,
			"created_on":  group.CreatedOn,
			"updated_on":  group.UpdatedOn,
			"workspace":   group.Workspace,
			"links":       group.Links,
		}
	}

	d.Set("groups", groups)
}
