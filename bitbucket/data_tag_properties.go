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

func dataTagProperties() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTagPropertiesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Required: true,
			},
			"properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataTagPropertiesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	tag := d.Get("tag").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataTagPropertiesRead", dumpResourceData(d, dataTagProperties().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/refs/tags/%s/properties", workspace, repoSlug, tag)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from tag properties call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate tag %s in repository %s/%s", tag, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading tag properties with params (%s): ", dumpResourceData(d, dataTagProperties().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	propertiesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] tag properties response: %v", propertiesBody)

	var propertiesResponse TagPropertiesResponse
	decodeerr := json.Unmarshal(propertiesBody, &propertiesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/refs/tags/%s/properties", workspace, repoSlug, tag))
	flattenTagProperties(&propertiesResponse, d)
	return nil
}

// TagPropertiesResponse represents the response from the tag properties API
type TagPropertiesResponse struct {
	Values []TagProperty `json:"values"`
	Page   int           `json:"page"`
	Size   int           `json:"size"`
	Next   string        `json:"next"`
}

// TagProperty represents a tag property
type TagProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Flattens the tag properties information
func flattenTagProperties(c *TagPropertiesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	properties := make([]interface{}, len(c.Values))
	for i, prop := range c.Values {
		properties[i] = map[string]interface{}{
			"key":   prop.Key,
			"value": prop.Value,
		}
	}

	d.Set("properties", properties)
}
