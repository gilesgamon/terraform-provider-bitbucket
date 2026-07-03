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

func dataPullRequestConflicts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPullRequestConflictsRead,
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
				Type:     schema.TypeString,
				Required: true,
			},
			"conflicts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: fileConflictSchema(),
				},
			},
		},
	}
}

func dataPullRequestConflictsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pullRequestID := d.Get("pull_request_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPullRequestConflictsRead", dumpResourceData(d, dataPullRequestConflicts().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pullrequests/%s/conflicts", workspace, repoSlug, pullRequestID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pull request conflicts call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pull request %s in repository %s/%s", pullRequestID, workspace, repoSlug)
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	conflictsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)

	var conflictsResponse FileConflictsResponse
	decodeerr := json.Unmarshal(conflictsBody, &conflictsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pullrequests/%s/conflicts", workspace, repoSlug, pullRequestID))
	d.Set("conflicts", flattenFileConflicts(conflictsResponse.Values))
	return nil
}
