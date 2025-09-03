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

func dataRepositoryFiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryFilesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ref": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Branch, tag, or commit hash",
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to directory or file",
			},
			"files": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"hash": {
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

func dataRepositoryFilesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	ref := d.Get("ref").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryFilesRead", dumpResourceData(d, dataRepositoryFiles().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/src/%s", workspace, repoSlug, ref)

	// Add path parameter if specified
	if path, ok := d.GetOk("path"); ok {
		url += "/" + path.(string)
	}

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository files call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s or ref %s", workspace, repoSlug, ref)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository files with params (%s): ", dumpResourceData(d, dataRepositoryFiles().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	filesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository files response: %v", filesBody)

	var filesResponse RepositoryFilesResponse
	decodeerr := json.Unmarshal(filesBody, &filesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/src/%s", workspace, repoSlug, ref))
	flattenRepositoryFiles(&filesResponse, d)
	return nil
}

// RepositoryFilesResponse represents the response from the repository files API
type RepositoryFilesResponse struct {
	Values []RepositoryFile `json:"values"`
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Next   string           `json:"next"`
}

// RepositoryFile represents a file in a repository
type RepositoryFile struct {
	Path  string                 `json:"path"`
	Type  string                 `json:"type"`
	Size  int                    `json:"size"`
	Hash  string                 `json:"hash"`
	Links map[string]interface{} `json:"links"`
}

// Flattens the repository files information
func flattenRepositoryFiles(c *RepositoryFilesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	files := make([]interface{}, len(c.Values))
	for i, file := range c.Values {
		files[i] = map[string]interface{}{
			"path":  file.Path,
			"type":  file.Type,
			"size":  file.Size,
			"hash":  file.Hash,
			"links": file.Links,
		}
	}

	d.Set("files", files)
}
