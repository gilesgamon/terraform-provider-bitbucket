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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The first revision (branch name, tag or commit hash).",
			},
			"target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The second revision (branch name, tag or commit hash).",
			},
			"hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hash of the common ancestor commit.",
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date": {
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
			"parents": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
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

	// The merge-base endpoint takes a single revspec containing exactly two
	// revisions separated by two dots.
	revspec := fmt.Sprintf("%s..%s", source, target)
	url := fmt.Sprintf("2.0/repositories/%s/%s/merge-base/%s", workspace, repoSlug, revspec)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from merge base call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate a common ancestor for %s and %s in repository %s/%s", source, target, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading merge base with params (%s): ", dumpResourceData(d, dataBranchMergeBase().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	mergeBaseBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)

	var commit Commit
	decodeerr := json.Unmarshal(mergeBaseBody, &commit)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/merge-base/%s", workspace, repoSlug, revspec))
	flattenCommit(&commit, d)
	return nil
}
