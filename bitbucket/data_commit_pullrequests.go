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

func dataCommitPullrequests() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCommitPullrequestsRead,
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
			"pullrequests": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"author": {
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
						"destination": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_on": {
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

func dataCommitPullrequestsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	commit := d.Get("commit").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitPullrequestsRead", dumpResourceData(d, dataCommitPullrequests().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commits/%s/pullrequests", workspace, repoSlug, commit)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from commit pull requests call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commit, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading commit pull requests with params (%s): ", dumpResourceData(d, dataCommitPullrequests().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	pullrequestsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] commit pull requests response: %v", pullrequestsBody)

	var pullrequestsResponse CommitPullrequestsResponse
	decodeerr := json.Unmarshal(pullrequestsBody, &pullrequestsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/commits/%s/pullrequests", workspace, repoSlug, commit))
	flattenCommitPullrequests(&pullrequestsResponse, d)
	return nil
}

// CommitPullrequestsResponse represents the response from the commit pull requests API
type CommitPullrequestsResponse struct {
	Values []CommitPullrequest `json:"values"`
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	Next   string              `json:"next"`
}

// CommitPullrequest represents a pull request containing a commit
type CommitPullrequest struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	State       string                 `json:"state"`
	Author      map[string]interface{} `json:"author"`
	Source      map[string]interface{} `json:"source"`
	Destination map[string]interface{} `json:"destination"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the commit pull requests information
func flattenCommitPullrequests(c *CommitPullrequestsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	pullrequests := make([]interface{}, len(c.Values))
	for i, pr := range c.Values {
		pullrequests[i] = map[string]interface{}{
			"id":           pr.ID,
			"title":        pr.Title,
			"description":  pr.Description,
			"state":        pr.State,
			"author":       pr.Author,
			"source":       pr.Source,
			"destination":  pr.Destination,
			"created_on":   pr.CreatedOn,
			"updated_on":   pr.UpdatedOn,
			"links":        pr.Links,
		}
	}

	d.Set("pullrequests", pullrequests)
}
