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

func dataRepositoryRefs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryRefsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search query string for ref names",
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort order (name, -name, target, -target)",
			},
			"refs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target": {
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

func dataRepositoryRefsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryRefsRead", dumpResourceData(d, dataRepositoryRefs().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/refs", workspace, repoSlug)

	// Build query parameters
	params := make(map[string]string)
	if q, ok := d.GetOk("q"); ok {
		params["q"] = q.(string)
	}
	if sort, ok := d.GetOk("sort"); ok {
		params["sort"] = sort.(string)
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
		return diag.Errorf("no response returned from repository refs call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository refs with params (%s): ", dumpResourceData(d, dataRepositoryRefs().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	refsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository refs response: %v", refsBody)

	var refsResponse RepositoryRefsResponse
	decodeerr := json.Unmarshal(refsBody, &refsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/refs", workspace, repoSlug))
	flattenRepositoryRefs(&refsResponse, d)
	return nil
}

// RepositoryRefsResponse represents the response from the repository refs API
type RepositoryRefsResponse struct {
	Values []RepositoryRef `json:"values"`
	Page   int             `json:"page"`
	Size   int             `json:"size"`
	Next   string          `json:"next"`
}

// RepositoryRef represents a reference in a repository
type RepositoryRef struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Target map[string]interface{} `json:"target"`
	Links  map[string]interface{} `json:"links"`
}

// Flattens the repository refs information
func flattenRepositoryRefs(c *RepositoryRefsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	refs := make([]interface{}, len(c.Values))
	for i, ref := range c.Values {
		refs[i] = map[string]interface{}{
			"name":   ref.Name,
			"type":   ref.Type,
			"target": ref.Target,
			"links":  ref.Links,
		}
	}

	d.Set("refs", refs)
}
