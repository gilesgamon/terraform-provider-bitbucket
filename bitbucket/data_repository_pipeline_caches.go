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

func dataRepositoryPipelineCaches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryPipelineCachesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"caches": {
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
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_accessed": {
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

func dataRepositoryPipelineCachesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryPipelineCachesRead", dumpResourceData(d, dataRepositoryPipelineCaches().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/caches", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository pipeline caches call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository pipeline caches with params (%s): ", dumpResourceData(d, dataRepositoryPipelineCaches().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	cachesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository pipeline caches response: %v", cachesBody)

	var cachesResponse RepositoryPipelineCachesResponse
	decodeerr := json.Unmarshal(cachesBody, &cachesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines_config/caches", workspace, repoSlug))
	flattenRepositoryPipelineCaches(&cachesResponse, d)
	return nil
}

// RepositoryPipelineCachesResponse represents the response from the repository pipeline caches API
type RepositoryPipelineCachesResponse struct {
	Values []RepositoryPipelineCache `json:"values"`
	Page   int                       `json:"page"`
	Size   int                       `json:"size"`
	Next   string                    `json:"next"`
}

// RepositoryPipelineCache represents a pipeline cache
type RepositoryPipelineCache struct {
	UUID         string                 `json:"uuid"`
	Name         string                 `json:"name"`
	Path         string                 `json:"path"`
	Size         int                    `json:"size"`
	LastAccessed string                 `json:"last_accessed"`
	CreatedOn    string                 `json:"created_on"`
	UpdatedOn    string                 `json:"updated_on"`
	Links        map[string]interface{} `json:"links"`
}

// Flattens the repository pipeline caches information
func flattenRepositoryPipelineCaches(c *RepositoryPipelineCachesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	caches := make([]interface{}, len(c.Values))
	for i, cache := range c.Values {
		caches[i] = map[string]interface{}{
			"uuid":          cache.UUID,
			"name":          cache.Name,
			"path":          cache.Path,
			"size":          cache.Size,
			"last_accessed": cache.LastAccessed,
			"created_on":    cache.CreatedOn,
			"updated_on":    cache.UpdatedOn,
			"links":         cache.Links,
		}
	}

	d.Set("caches", caches)
}
