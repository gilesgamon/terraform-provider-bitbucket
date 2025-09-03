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

func dataBranchMergeBase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataBranchMergeBaseRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target": {
				Type:     schema.TypeString,
				Required: true,
			},
			"merge_base": {
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
	}
}

func dataBranchMergeBaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	source := d.Get("source").(string)
	target := d.Get("target").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataBranchMergeBaseRead", dumpResourceData(d, dataBranchMergeBase().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commits/%s/merge-base/%s", workspace, repoSlug, source, target)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from branch merge base call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate branches %s or %s in repository %s/%s", source, target, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading branch merge base with params (%s): ", dumpResourceData(d, dataBranchMergeBase().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	mergeBaseBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] branch merge base response: %v", mergeBaseBody)

	var mergeBase BranchMergeBaseData
	decodeerr := json.Unmarshal(mergeBaseBody, &mergeBase)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/commits/%s/merge-base/%s", workspace, repoSlug, source, target))
	flattenBranchMergeBase(&mergeBase, d)
	return nil
}

// BranchMergeBaseData represents the merge base between two branches
type BranchMergeBaseData struct {
	MergeBase string                 `json:"merge_base"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the branch merge base information
func flattenBranchMergeBase(c *BranchMergeBaseData, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("merge_base", c.MergeBase)
	d.Set("links", c.Links)
}
