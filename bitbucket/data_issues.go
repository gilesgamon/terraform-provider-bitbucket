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

func dataIssues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssuesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issues": {
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
						"content": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assignee": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"reporter": {
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

func dataIssuesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssuesRead", dumpResourceData(d, dataIssues().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issues call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issues with params (%s): ", dumpResourceData(d, dataIssues().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	issuesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issues response: %v", issuesBody)

	var issuesResponse IssuesResponse
	decodeerr := json.Unmarshal(issuesBody, &issuesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues", workspace, repoSlug))
	flattenIssues(&issuesResponse, d)
	return nil
}

// IssuesResponse represents the response from the issues API
type IssuesResponse struct {
	Values []IssueData `json:"values"`
	Page   int         `json:"page"`
	Size   int         `json:"size"`
	Next   string      `json:"next"`
}

// IssueData represents an issue in the list
type IssueData struct {
	ID        int                    `json:"id"`
	Title     string                 `json:"title"`
	Content   map[string]interface{} `json:"content"`
	State     string                 `json:"state"`
	Kind      string                 `json:"kind"`
	Priority  string                 `json:"priority"`
	Assignee  map[string]interface{} `json:"assignee"`
	Reporter  map[string]interface{} `json:"reporter"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issues information
func flattenIssues(c *IssuesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	issues := make([]interface{}, len(c.Values))
	for i, issue := range c.Values {
		issues[i] = map[string]interface{}{
			"id":         issue.ID,
			"title":      issue.Title,
			"content":    issue.Content,
			"state":      issue.State,
			"kind":       issue.Kind,
			"priority":   issue.Priority,
			"assignee":   issue.Assignee,
			"reporter":   issue.Reporter,
			"created_on": issue.CreatedOn,
			"updated_on": issue.UpdatedOn,
			"links":      issue.Links,
		}
	}

	d.Set("issues", issues)
}
