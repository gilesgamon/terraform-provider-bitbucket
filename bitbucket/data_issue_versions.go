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

func dataIssueVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueVersionsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"versions": {
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
						"released": {
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

func dataIssueVersionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueVersionsRead", dumpResourceData(d, dataIssueVersions().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/versions", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue versions call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue versions with params (%s): ", dumpResourceData(d, dataIssueVersions().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	versionsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue versions response: %v", versionsBody)

	var versionsResponse IssueVersionsResponse
	decodeerr := json.Unmarshal(versionsBody, &versionsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/versions", workspace, repoSlug))
	flattenIssueVersions(&versionsResponse, d)
	return nil
}

// IssueVersionsResponse represents the response from the issue versions API
type IssueVersionsResponse struct {
	Values []IssueVersion `json:"values"`
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Next   string         `json:"next"`
}

// IssueVersion represents an issue version
type IssueVersion struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Released    bool                   `json:"released"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the issue versions information
func flattenIssueVersions(c *IssueVersionsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	versions := make([]interface{}, len(c.Values))
	for i, version := range c.Values {
		versions[i] = map[string]interface{}{
			"uuid":        version.UUID,
			"name":        version.Name,
			"description": version.Description,
			"released":    version.Released,
			"created_on":  version.CreatedOn,
			"updated_on":  version.UpdatedOn,
			"links":       version.Links,
		}
	}

	d.Set("versions", versions)
}
