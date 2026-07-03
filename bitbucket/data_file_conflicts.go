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

func dataFileConflicts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFileConflictsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"spec": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A merge base revspec, e.g. `main..feature-branch`, used to compute file conflicts.",
			},
			"conflicts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: fileConflictSchema(),
				},
			},
		},
	}
}

func dataFileConflictsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	spec := d.Get("spec").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataFileConflictsRead", dumpResourceData(d, dataFileConflicts().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/file-conflicts/%s", workspace, repoSlug, spec)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from file conflicts call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s or spec %s", workspace, repoSlug, spec)
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	conflictsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)

	var conflictsResponse FileConflictsResponse
	decodeerr := json.Unmarshal(conflictsBody, &conflictsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/file-conflicts/%s", workspace, repoSlug, spec))
	d.Set("conflicts", flattenFileConflicts(conflictsResponse.Values))
	return nil
}

// FileConflictsResponse represents a paginated list of file conflicts
type FileConflictsResponse struct {
	Values  []FileConflict `json:"values"`
	Page    int            `json:"page"`
	Size    int            `json:"size"`
	Pagelen int            `json:"pagelen"`
	Next    string         `json:"next"`
}

// FileConflict represents a single file conflict object
type FileConflict struct {
	Type     string `json:"type"`
	Path     string `json:"path"`
	Scenario string `json:"scenario"`
	Message  string `json:"message"`
}

// fileConflictSchema returns the shared schema describing a single file conflict
func fileConflictSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"path": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"scenario": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"message": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// flattenFileConflicts converts file conflicts into a Terraform-friendly structure
func flattenFileConflicts(conflicts []FileConflict) []interface{} {
	result := make([]interface{}, len(conflicts))
	for i, conflict := range conflicts {
		result[i] = map[string]interface{}{
			"type":     conflict.Type,
			"path":     conflict.Path,
			"scenario": conflict.Scenario,
			"message":  conflict.Message,
		}
	}
	return result
}
