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

func dataRepositoryHooks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryHooksRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
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
						"skip_cert_verification": {
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
					},
				},
			},
		},
	}
}

func dataRepositoryHooksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryHooksRead", dumpResourceData(d, dataRepositoryHooks().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/hooks", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository hooks call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository hooks with params (%s): ", dumpResourceData(d, dataRepositoryHooks().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	hooksBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository hooks response: %v", hooksBody)

	var hooksResponse RepositoryHooksResponse
	decodeerr := json.Unmarshal(hooksBody, &hooksResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/hooks", workspace, repoSlug))
	flattenRepositoryHooks(&hooksResponse, d)
	return nil
}

// RepositoryHooksResponse represents the response from the repository hooks API
type RepositoryHooksResponse struct {
	Values []RepositoryHook `json:"values"`
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Next   string           `json:"next"`
}

// RepositoryHook represents a webhook in a repository
type RepositoryHook struct {
	UUID                 string   `json:"uuid"`
	Description          string   `json:"description"`
	URL                  string   `json:"url"`
	Active               bool     `json:"active"`
	Events               []string `json:"events"`
	SkipCertVerification bool     `json:"skip_cert_verification"`
	CreatedOn            string   `json:"created_on"`
	UpdatedOn            string   `json:"updated_on"`
}

// Flattens the repository hooks information
func flattenRepositoryHooks(c *RepositoryHooksResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	hooks := make([]interface{}, len(c.Values))
	for i, hook := range c.Values {
		hooks[i] = map[string]interface{}{
			"uuid":                   hook.UUID,
			"description":            hook.Description,
			"url":                    hook.URL,
			"active":                 hook.Active,
			"events":                 hook.Events,
			"skip_cert_verification": hook.SkipCertVerification,
			"created_on":             hook.CreatedOn,
			"updated_on":             hook.UpdatedOn,
		}
	}

	d.Set("hooks", hooks)
}
