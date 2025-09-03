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

func dataIssueTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueTypesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"types": {
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

func dataIssueTypesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueTypesRead", dumpResourceData(d, dataIssueTypes().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-types", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue types call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue types with params (%s): ", dumpResourceData(d, dataIssueTypes().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	typesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue types response: %v", typesBody)

	var typesResponse IssueTypesResponse
	decodeerr := json.Unmarshal(typesBody, &typesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-types", workspace, repoSlug))
	flattenIssueTypes(&typesResponse, d)
	return nil
}

// IssueTypesResponse represents the response from the issue types API
type IssueTypesResponse struct {
	Values []IssueType `json:"values"`
	Page   int         `json:"page"`
	Size   int         `json:"size"`
	Next   string      `json:"next"`
}

// IssueType represents an issue type
type IssueType struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Icon        string                 `json:"icon"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the issue types information
func flattenIssueTypes(c *IssueTypesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	types := make([]interface{}, len(c.Values))
	for i, issueType := range c.Values {
		types[i] = map[string]interface{}{
			"uuid":        issueType.UUID,
			"name":        issueType.Name,
			"description": issueType.Description,
			"icon":        issueType.Icon,
			"created_on":  issueType.CreatedOn,
			"updated_on":  issueType.UpdatedOn,
			"links":       issueType.Links,
		}
	}

	d.Set("types", types)
}
