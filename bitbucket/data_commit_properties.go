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

func dataCommitProperties() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCommitPropertiesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commit": {
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

func dataCommitPropertiesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	commit := d.Get("commit").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitPropertiesRead", dumpResourceData(d, dataCommitProperties().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commits/%s/properties", workspace, repoSlug, commit)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from commit properties call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commit, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading commit properties with params (%s): ", dumpResourceData(d, dataCommitProperties().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	propertiesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] commit properties response: %v", propertiesBody)

	var propertiesResponse CommitPropertiesResponse
	decodeerr := json.Unmarshal(propertiesBody, &propertiesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/commits/%s/properties", workspace, repoSlug, commit))
	flattenCommitProperties(&propertiesResponse, d)
	return nil
}

// CommitPropertiesResponse represents the response from the commit properties API
type CommitPropertiesResponse struct {
	Values []CommitProperty `json:"values"`
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Next   string           `json:"next"`
}

// CommitProperty represents a commit property
type CommitProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Flattens the commit properties information
func flattenCommitProperties(c *CommitPropertiesResponse, d *schema.ResourceData) {
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
