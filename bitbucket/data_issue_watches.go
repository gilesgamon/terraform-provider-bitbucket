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

func dataIssueWatches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueWatchesRead,
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
				Type:     schema.TypeString,
				Required: true,
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

func dataIssueWatchesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	issueID := d.Get("issue_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueWatchesRead", dumpResourceData(d, dataIssueWatches().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues/%s/watches", workspace, repoSlug, issueID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue watches call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue %s in repository %s/%s", issueID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue watches with params (%s): ", dumpResourceData(d, dataIssueWatches().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	watchesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue watches response: %v", watchesBody)

	var watchesResponse IssueWatchesResponse
	decodeerr := json.Unmarshal(watchesBody, &watchesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues/%s/watches", workspace, repoSlug, issueID))
	flattenIssueWatches(&watchesResponse, d)
	return nil
}

// IssueWatchesResponse represents the response from the issue watches API
type IssueWatchesResponse struct {
	Values []IssueWatch `json:"values"`
	Page   int          `json:"page"`
	Size   int          `json:"size"`
	Next   string       `json:"next"`
}

// IssueWatch represents a watch on an issue
type IssueWatch struct {
	User      map[string]interface{} `json:"user"`
	CreatedOn string                 `json:"created_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue watches information
func flattenIssueWatches(c *IssueWatchesResponse, d *schema.ResourceData) {
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
