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

func dataRepositoryIssueWatches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryIssueWatchesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issue_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Issue ID",
			},
			"watches": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
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

func dataRepositoryIssueWatchesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	issueID := d.Get("issue_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryIssueWatchesRead", dumpResourceData(d, dataRepositoryIssueWatches().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues/%s/watches", workspace, repoSlug, issueID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository issue watches call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue %s in repository %s/%s", issueID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository issue watches with params (%s): ", dumpResourceData(d, dataRepositoryIssueWatches().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	watchesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository issue watches response: %v", watchesBody)

	var watchesResponse RepositoryIssueWatchesResponse
	decodeerr := json.Unmarshal(watchesBody, &watchesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues/%s/watches", workspace, repoSlug, issueID))
	flattenRepositoryIssueWatches(&watchesResponse, d)
	return nil
}

// RepositoryIssueWatchesResponse represents the response from the repository issue watches API
type RepositoryIssueWatchesResponse struct {
	Values []RepositoryIssueWatch `json:"values"`
	Page   int                    `json:"page"`
	Size   int                    `json:"size"`
	Next   string                 `json:"next"`
}

// RepositoryIssueWatch represents an issue watch
type RepositoryIssueWatch struct {
	User      map[string]interface{} `json:"user"`
	CreatedOn string                 `json:"created_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the repository issue watches information
func flattenRepositoryIssueWatches(c *RepositoryIssueWatchesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	watches := make([]interface{}, len(c.Values))
	for i, watch := range c.Values {
		watches[i] = map[string]interface{}{
			"user":       watch.User,
			"created_on": watch.CreatedOn,
			"links":      watch.Links,
		}
	}

	d.Set("watches", watches)
}
