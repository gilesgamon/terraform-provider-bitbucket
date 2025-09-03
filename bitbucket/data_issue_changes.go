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

func dataIssueChanges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueChangesRead,
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
			"changes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
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
						"changes": {
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

func dataIssueChangesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	issueID := d.Get("issue_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueChangesRead", dumpResourceData(d, dataIssueChanges().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues/%s/changes", workspace, repoSlug, issueID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue changes call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue %s in repository %s/%s", issueID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue changes with params (%s): ", dumpResourceData(d, dataIssueChanges().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	changesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue changes response: %v", changesBody)

	var changesResponse IssueChangesResponse
	decodeerr := json.Unmarshal(changesBody, &changesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues/%s/changes", workspace, repoSlug, issueID))
	flattenIssueChanges(&changesResponse, d)
	return nil
}

// IssueChangesResponse represents the response from the issue changes API
type IssueChangesResponse struct {
	Values []IssueChange `json:"values"`
	Page   int           `json:"page"`
	Size   int           `json:"size"`
	Next   string        `json:"next"`
}

// IssueChange represents a change to an issue
type IssueChange struct {
	UUID      string                 `json:"uuid"`
	User      map[string]interface{} `json:"user"`
	Changes   map[string]interface{} `json:"changes"`
	CreatedOn string                 `json:"created_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue changes information
func flattenIssueChanges(c *IssueChangesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	changes := make([]interface{}, len(c.Values))
	for i, change := range c.Values {
		changes[i] = map[string]interface{}{
			"uuid":       change.UUID,
			"user":       change.User,
			"changes":    change.Changes,
			"created_on": change.CreatedOn,
			"links":      change.Links,
		}
	}

	d.Set("changes", changes)
}
