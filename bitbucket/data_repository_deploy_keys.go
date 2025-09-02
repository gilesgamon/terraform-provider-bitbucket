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

func dataRepositoryDeployKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryDeployKeysRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deploy_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_used": {
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

func dataRepositoryDeployKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryDeployKeysRead", dumpResourceData(d, dataRepositoryDeployKeys().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/deploy-keys", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository deploy keys call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository deploy keys with params (%s): ", dumpResourceData(d, dataRepositoryDeployKeys().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	deployKeysBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository deploy keys response: %v", deployKeysBody)

	var deployKeysResponse RepositoryDeployKeysResponse
	decodeerr := json.Unmarshal(deployKeysBody, &deployKeysResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/deploy-keys", workspace, repoSlug))
	flattenRepositoryDeployKeys(&deployKeysResponse, d)
	return nil
}

// RepositoryDeployKeysResponse represents the response from the repository deploy keys API
type RepositoryDeployKeysResponse struct {
	Values []RepositoryDeployKey `json:"values"`
	Page   int                   `json:"page"`
	Size   int                   `json:"size"`
	Next   string                `json:"next"`
}

// RepositoryDeployKey represents a deploy key in a repository
type RepositoryDeployKey struct {
	ID        int                    `json:"id"`
	Key       string                 `json:"key"`
	Label     string                 `json:"label"`
	Comment   string                 `json:"comment"`
	CreatedOn string                 `json:"created_on"`
	LastUsed  string                 `json:"last_used"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the repository deploy keys information
func flattenRepositoryDeployKeys(c *RepositoryDeployKeysResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	deployKeys := make([]interface{}, len(c.Values))
	for i, key := range c.Values {
		deployKeys[i] = map[string]interface{}{
			"id":         key.ID,
			"key":        key.Key,
			"label":      key.Label,
			"comment":    key.Comment,
			"created_on": key.CreatedOn,
			"last_used":  key.LastUsed,
			"links":      key.Links,
		}
	}

	d.Set("deploy_keys", deployKeys)
}
