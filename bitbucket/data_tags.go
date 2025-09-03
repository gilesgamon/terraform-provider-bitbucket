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

func dataTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTagsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tagger": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"date": {
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

func dataTagsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataTagsRead", dumpResourceData(d, dataTags().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/refs/tags", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from tags call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading tags with params (%s): ", dumpResourceData(d, dataTags().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	tagsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] tags response: %v", tagsBody)

	var tagsResponse TagsResponse
	decodeerr := json.Unmarshal(tagsBody, &tagsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/refs/tags", workspace, repoSlug))
	flattenTags(&tagsResponse, d)
	return nil
}

// TagsResponse represents the response from the tags API
type TagsResponse struct {
	Values []TagData `json:"values"`
	Page   int       `json:"page"`
	Size   int       `json:"size"`
	Next   string    `json:"next"`
}

// TagData represents a tag in the list
type TagData struct {
	Name   string                 `json:"name"`
	Target map[string]interface{} `json:"target"`
	Message string                 `json:"message"`
	Tagger map[string]interface{} `json:"tagger"`
	Date   string                 `json:"date"`
	Links  map[string]interface{} `json:"links"`
}

// Flattens the tags information
func flattenTags(c *TagsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	tags := make([]interface{}, len(c.Values))
	for i, tag := range c.Values {
		tags[i] = map[string]interface{}{
			"name":   tag.Name,
			"target": tag.Target,
			"message": tag.Message,
			"tagger": tag.Tagger,
			"date":   tag.Date,
			"links":  tag.Links,
		}
	}

	d.Set("tags", tags)
}
