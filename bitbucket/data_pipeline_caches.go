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

func dataPipelineCaches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineCachesRead,
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
					},
				},
			},
		},
	}
}

func dataPipelineCachesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineCachesRead", dumpResourceData(d, dataPipelineCaches().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/caches", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline caches call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline caches for repository %s", repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline caches with params (%s): ", dumpResourceData(d, dataPipelineCaches().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	cachesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline caches response: %v", cachesBody)

	var cachesResponse PipelineCachesResponse
	decodeerr := json.Unmarshal(cachesBody, &cachesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines_config/caches", workspace, repoSlug))
	flattenPipelineCaches(&cachesResponse, d)
	return nil
}

// PipelineCachesResponse represents the response from the pipeline caches API
type PipelineCachesResponse struct {
	Values []PipelineCache `json:"values"`
}

// PipelineCache represents a pipeline cache
type PipelineCache struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	Size         int    `json:"size"`
	LastAccessed string `json:"last_accessed"`
	CreatedOn    string `json:"created_on"`
}

// Flattens the pipeline caches information
func flattenPipelineCaches(c *PipelineCachesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	caches := make([]interface{}, len(c.Values))
	for i, cache := range c.Values {
		caches[i] = map[string]interface{}{
			"name":          cache.Name,
			"path":          cache.Path,
			"size":          cache.Size,
			"last_accessed": cache.LastAccessed,
			"created_on":    cache.CreatedOn,
		}
	}

	d.Set("caches", caches)
}
