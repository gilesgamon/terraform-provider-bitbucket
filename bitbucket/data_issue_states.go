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

func dataIssueStates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueStatesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"states": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"icon": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"color": {
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

func dataIssueStatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueStatesRead", dumpResourceData(d, dataIssueStates().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-states", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue states call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue states with params (%s): ", dumpResourceData(d, dataIssueStates().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	statesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue states response: %v", statesBody)

	var statesResponse IssueStatesResponse
	decodeerr := json.Unmarshal(statesBody, &statesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-states", workspace, repoSlug))
	flattenIssueStates(&statesResponse, d)
	return nil
}

// IssueStatesResponse represents the response from the issue states API
type IssueStatesResponse struct {
	Values []IssueState `json:"values"`
	Page   int          `json:"page"`
	Size   int          `json:"size"`
	Next   string       `json:"next"`
}

// IssueState represents an issue state
type IssueState struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Icon        string                 `json:"icon"`
	Color       string                 `json:"color"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the issue states information
func flattenIssueStates(c *IssueStatesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	states := make([]interface{}, len(c.Values))
	for i, state := range c.Values {
		states[i] = map[string]interface{}{
			"uuid":        state.UUID,
			"name":        state.Name,
			"description": state.Description,
			"icon":        state.Icon,
			"color":       state.Color,
			"created_on":  state.CreatedOn,
			"updated_on":  state.UpdatedOn,
			"links":       state.Links,
		}
	}

	d.Set("states", states)
}
