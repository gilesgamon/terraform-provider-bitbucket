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

func dataRepositoryDownloads() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryDownloadsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"downloads": {
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
						"downloads": {
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

func dataRepositoryDownloadsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryDownloadsRead", dumpResourceData(d, dataRepositoryDownloads().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/downloads", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository downloads call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository downloads with params (%s): ", dumpResourceData(d, dataRepositoryDownloads().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	downloadsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository downloads response: %v", downloadsBody)

	var downloadsResponse RepositoryDownloadsResponse
	decodeerr := json.Unmarshal(downloadsBody, &downloadsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/downloads", workspace, repoSlug))
	flattenRepositoryDownloads(&downloadsResponse, d)
	return nil
}

// RepositoryDownloadsResponse represents the response from the repository downloads API
type RepositoryDownloadsResponse struct {
	Values []RepositoryDownload `json:"values"`
	Page   int                  `json:"page"`
	Size   int                  `json:"size"`
	Next   string               `json:"next"`
}

// RepositoryDownload represents a download in a repository
type RepositoryDownload struct {
	Name       string                 `json:"name"`
	Size       int                    `json:"size"`
	Downloads  int                    `json:"downloads"`
	CreatedOn  string                 `json:"created_on"`
	Links      map[string]interface{} `json:"links"`
}

// Flattens the repository downloads information
func flattenRepositoryDownloads(c *RepositoryDownloadsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	downloads := make([]interface{}, len(c.Values))
	for i, download := range c.Values {
		downloads[i] = map[string]interface{}{
			"name":      download.Name,
			"size":      download.Size,
			"downloads": download.Downloads,
			"created_on": download.CreatedOn,
			"links":     download.Links,
		}
	}

	d.Set("downloads", downloads)
}
