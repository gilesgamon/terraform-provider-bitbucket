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

func dataRepositoryPatches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryPatchesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"patches": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
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

func dataRepositoryPatchesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryPatchesRead", dumpResourceData(d, dataRepositoryPatches().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/patches", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository patches call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository patches with params (%s): ", dumpResourceData(d, dataRepositoryPatches().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	patchesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository patches response: %v", patchesBody)

	var patchesResponse RepositoryPatchesResponse
	decodeerr := json.Unmarshal(patchesBody, &patchesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/patches", workspace, repoSlug))
	flattenRepositoryPatches(&patchesResponse, d)
	return nil
}

// RepositoryPatchesResponse represents the response from the repository patches API
type RepositoryPatchesResponse struct {
	Values []RepositoryPatch `json:"values"`
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Next   string            `json:"next"`
}

// RepositoryPatch represents a patch in a repository
type RepositoryPatch struct {
	Name      string                 `json:"name"`
	Size      int                    `json:"size"`
	CreatedOn string                 `json:"created_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the repository patches information
func flattenRepositoryPatches(c *RepositoryPatchesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	patches := make([]interface{}, len(c.Values))
	for i, patch := range c.Values {
		patches[i] = map[string]interface{}{
			"name":       patch.Name,
			"size":       patch.Size,
			"created_on": patch.CreatedOn,
			"links":      patch.Links,
		}
	}

	d.Set("patches", patches)
}
