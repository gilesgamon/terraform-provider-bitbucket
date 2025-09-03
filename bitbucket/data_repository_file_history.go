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

func dataRepositoryFileHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryFileHistoryRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
		"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to the file",
			},
			"revision": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Revision (commit hash) to get history for",
			},
			"history": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"commit": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
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

func dataRepositoryFileHistoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	path := d.Get("path").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryFileHistoryRead", dumpResourceData(d, dataRepositoryFileHistory().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/filehistory/%s", workspace, repoSlug, path)

	// Add revision parameter if specified
	if revision, ok := d.GetOk("revision"); ok {
		url += "?revision=" + revision.(string)
	}

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository file history call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate file %s in repository %s/%s", path, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository file history with params (%s): ", dumpResourceData(d, dataRepositoryFileHistory().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	historyBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository file history response: %v", historyBody)

	var historyResponse RepositoryFileHistoryResponse
	decodeerr := json.Unmarshal(historyBody, &historyResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/filehistory/%s", workspace, repoSlug, path))
	flattenRepositoryFileHistory(&historyResponse, d)
	return nil
}

// RepositoryFileHistoryResponse represents the response from the repository file history API
type RepositoryFileHistoryResponse struct {
	Values []RepositoryFileHistoryItem `json:"values"`
	Page   int                         `json:"page"`
	Size   int                         `json:"size"`
	Next   string                      `json:"next"`
}

// RepositoryFileHistoryItem represents a file history item
type RepositoryFileHistoryItem struct {
	Commit map[string]interface{} `json:"commit"`
	Path   string                 `json:"path"`
	Type   string                 `json:"type"`
	Size   int                    `json:"size"`
	Links  map[string]interface{} `json:"links"`
}

// Flattens the repository file history information
func flattenRepositoryFileHistory(c *RepositoryFileHistoryResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	history := make([]interface{}, len(c.Values))
	for i, item := range c.Values {
		history[i] = map[string]interface{}{
			"commit": item.Commit,
			"path":   item.Path,
			"type":   item.Type,
			"size":   item.Size,
			"links":  item.Links,
		}
	}

	d.Set("history", history)
}
