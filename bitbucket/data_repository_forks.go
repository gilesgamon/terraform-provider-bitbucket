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

func dataRepositoryForks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryForksRead,
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
				Description: "Search query string for repository names",
			},
			"forks": {
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
						"full_name": {
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
						"fork_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"language": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
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
						"project": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"mainbranch": {
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

func dataRepositoryForksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryForksRead", dumpResourceData(d, dataRepositoryForks().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/forks", workspace, repoSlug)

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
		return diag.Errorf("no response returned from repository forks call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository forks with params (%s): ", dumpResourceData(d, dataRepositoryForks().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	forksBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository forks response: %v", forksBody)

	var forksResponse RepositoryForksResponse
	decodeerr := json.Unmarshal(forksBody, &forksResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/forks", workspace, repoSlug))
	flattenRepositoryForks(&forksResponse, d)
	return nil
}

// RepositoryForksResponse represents the response from the repository forks API
type RepositoryForksResponse struct {
	Values []RepositoryFork `json:"values"`
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Next   string           `json:"next"`
}

// RepositoryFork represents a repository fork
type RepositoryFork struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	FullName    string                 `json:"full_name"`
	Description string                 `json:"description"`
	IsPrivate   bool                   `json:"is_private"`
	ForkPolicy  string                 `json:"fork_policy"`
	Language    string                 `json:"language"`
	Size        int                    `json:"size"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Workspace   map[string]interface{} `json:"workspace"`
	Project     map[string]interface{} `json:"project"`
	Mainbranch  map[string]interface{} `json:"mainbranch"`
}

// Flattens the repository forks information
func flattenRepositoryForks(c *RepositoryForksResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	forks := make([]interface{}, len(c.Values))
	for i, fork := range c.Values {
		forks[i] = map[string]interface{}{
			"uuid":        fork.UUID,
			"name":        fork.Name,
			"full_name":   fork.FullName,
			"description": fork.Description,
			"is_private":  fork.IsPrivate,
			"fork_policy": fork.ForkPolicy,
			"language":    fork.Language,
			"size":        fork.Size,
			"created_on":  fork.CreatedOn,
			"updated_on":  fork.UpdatedOn,
			"workspace":   fork.Workspace,
			"project":     fork.Project,
			"mainbranch":  fork.Mainbranch,
		}
	}

	d.Set("forks", forks)
}
