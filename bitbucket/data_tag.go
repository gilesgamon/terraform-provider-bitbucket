package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/DrFaust92/bitbucket-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataTag() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTagRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"author": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	tagName := d.Get("tag_name").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataTagRead", dumpResourceData(d, dataTag().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/refs/tags/%s",
		workspace,
		repoSlug,
		tagName,
	)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from tag call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate tag %s in repository %s/%s", tagName, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading tag information with params (%s): ", dumpResourceData(d, dataTag().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	tagBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] tag response: %v", tagBody)

	var tag Tag
	decodeerr := json.Unmarshal(tagBody, &tag)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", workspace, repoSlug, tag.Name))
	flattenTag(&tag, d)
	return nil
}

// Tag represents a Bitbucket tag
type Tag struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Target     TagTarget              `json:"target"`
	Hash       string                 `json:"hash"`
	Repository bitbucket.Repository   `json:"repository"`
	Links      map[string]interface{} `json:"links"`
}

// TagTarget represents the target of a tag (usually a commit)
type TagTarget struct {
	Hash  string                 `json:"hash"`
	Type  string                 `json:"type"`
	Links map[string]interface{} `json:"links"`
}

// Flattens the tag information
func flattenTag(t *Tag, d *schema.ResourceData) {
	if t == nil {
		return
	}

	d.Set("uuid", t.Name)
	d.Set("name", t.Name)
	d.Set("target_hash", t.Target.Hash)
	d.Set("target_date", t.Target.Type)
	d.Set("message", t.Type)
}
