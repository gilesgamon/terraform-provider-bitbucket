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

func dataRepositoryIssueChanges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryIssueChangesRead,
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
			"changes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
						"changes": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func dataRepositoryIssueChangesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	issueID := d.Get("issue_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryIssueChangesRead", dumpResourceData(d, dataRepositoryIssueChanges().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues/%s/changes", workspace, repoSlug, issueID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository issue changes call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue %s in repository %s/%s", issueID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository issue changes with params (%s): ", dumpResourceData(d, dataRepositoryIssueChanges().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	changesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository issue changes response: %v", changesBody)

	var changesResponse RepositoryIssueChangesResponse
	decodeerr := json.Unmarshal(changesBody, &changesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues/%s/changes", workspace, repoSlug, issueID))
	flattenRepositoryIssueChanges(&changesResponse, d)
	return nil
}

// RepositoryIssueChangesResponse represents the response from the repository issue changes API
type RepositoryIssueChangesResponse struct {
	Values []RepositoryIssueChange `json:"values"`
	Page   int                     `json:"page"`
	Size   int                     `json:"size"`
	Next   string                  `json:"next"`
}

// RepositoryIssueChange represents an issue change
type RepositoryIssueChange struct {
	ID        int                    `json:"id"`
	Type      string                 `json:"type"`
	User      map[string]interface{} `json:"user"`
	CreatedOn string                 `json:"created_on"`
	Changes   map[string]interface{} `json:"changes"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the repository issue changes information
func flattenRepositoryIssueChanges(c *RepositoryIssueChangesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	changes := make([]interface{}, len(c.Values))
	for i, change := range c.Values {
		changes[i] = map[string]interface{}{
			"id":         change.ID,
			"type":       change.Type,
			"user":       change.User,
			"created_on": change.CreatedOn,
			"changes":    change.Changes,
			"links":      change.Links,
		}
	}

	d.Set("changes", changes)
}
