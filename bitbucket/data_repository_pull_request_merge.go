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

func dataRepositoryPullRequestMerge() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryPullRequestMergeRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pull_request_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pull request ID",
			},
			"merge_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"merge_commit": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"close_source_branch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"merge_strategy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataRepositoryPullRequestMergeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pullRequestID := d.Get("pull_request_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryPullRequestMergeRead", dumpResourceData(d, dataRepositoryPullRequestMerge().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pullrequests/%s/merge", workspace, repoSlug, pullRequestID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository pull request merge call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pull request %s in repository %s/%s", pullRequestID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository pull request merge with params (%s): ", dumpResourceData(d, dataRepositoryPullRequestMerge().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	mergeBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository pull request merge response: %v", mergeBody)

	var mergeResponse RepositoryPullRequestMerge
	decodeerr := json.Unmarshal(mergeBody, &mergeResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pullrequests/%s/merge", workspace, repoSlug, pullRequestID))
	flattenRepositoryPullRequestMerge(&mergeResponse, d)
	return nil
}

// RepositoryPullRequestMerge represents the response from the repository pull request merge API
type RepositoryPullRequestMerge struct {
	MergeStatus        map[string]interface{} `json:"merge_status"`
	MergeCommit        map[string]interface{} `json:"merge_commit"`
	CloseSourceBranch  bool                   `json:"close_source_branch"`
	MergeStrategy      string                 `json:"merge_strategy"`
	Destination        map[string]interface{} `json:"destination"`
	Source             map[string]interface{} `json:"source"`
}

// Flattens the repository pull request merge information
func flattenRepositoryPullRequestMerge(c *RepositoryPullRequestMerge, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("merge_status", c.MergeStatus)
	d.Set("merge_commit", c.MergeCommit)
	d.Set("close_source_branch", c.CloseSourceBranch)
	d.Set("merge_strategy", c.MergeStrategy)
	d.Set("destination", c.Destination)
	d.Set("source", c.Source)
}
