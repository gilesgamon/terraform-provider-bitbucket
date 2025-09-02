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
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter issues by state (open, resolved, closed, declined, merged)",
			},
			"kind": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter issues by kind (bug, enhancement, proposal, task)",
			},
			"priority": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter issues by priority (trivial, minor, major, critical, blocker)",
			},
			"assignee": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter issues by assignee username",
			},
			"reporter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter issues by reporter username",
			},
			"milestone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter issues by milestone name",
			},
			"component": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter issues by component name",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter issues by version name",
			},
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search query string",
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort field (created_on, updated_on, priority, kind, state)",
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
							Type:     schema.TypeString,
							Computed: true,
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"reporter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"milestone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"votes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"watches": {
							Type:     schema.TypeInt,
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

	// Build query parameters
	params := make(map[string]string)
	if state, ok := d.GetOk("state"); ok {
		params["state"] = state.(string)
	}
	if kind, ok := d.GetOk("kind"); ok {
		params["kind"] = kind.(string)
	}
	if priority, ok := d.GetOk("priority"); ok {
		params["priority"] = priority.(string)
	}
	if assignee, ok := d.GetOk("assignee"); ok {
		params["assignee"] = assignee.(string)
	}
	if reporter, ok := d.GetOk("reporter"); ok {
		params["reporter"] = reporter.(string)
	}
	if milestone, ok := d.GetOk("milestone"); ok {
		params["milestone"] = milestone.(string)
	}
	if component, ok := d.GetOk("component"); ok {
		params["component"] = component.(string)
	}
	if version, ok := d.GetOk("version"); ok {
		params["version"] = version.(string)
	}
	if q, ok := d.GetOk("q"); ok {
		params["q"] = q.(string)
	}
	if sort, ok := d.GetOk("sort"); ok {
		params["sort"] = sort.(string)
	}

	// Add query parameters to URL
	if len(params) > 0 {
		url += "?"
		first := true
		for key, value := range params {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", key, value)
			first = false
		}
	}

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
	Values []Issue `json:"values"`
	Page   int     `json:"page"`
	Size   int     `json:"size"`
	Next   string  `json:"next"`
}

// Issue represents an issue in the repository
type Issue struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	State       string                 `json:"state"`
	Kind        string                 `json:"kind"`
	Priority    string                 `json:"priority"`
	Assignee    string                 `json:"assignee"`
	Reporter    string                 `json:"reporter"`
	Milestone   string                 `json:"milestone"`
	Component   string                 `json:"component"`
	Version     string                 `json:"version"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Votes       int                    `json:"votes"`
	Watches     int                    `json:"watches"`
	Links       map[string]interface{} `json:"links"`
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
			"milestone":  issue.Milestone,
			"component":  issue.Component,
			"version":    issue.Version,
			"created_on": issue.CreatedOn,
			"updated_on": issue.UpdatedOn,
			"votes":      issue.Votes,
			"watches":    issue.Watches,
			"links":      issue.Links,
		}
	}

	d.Set("issues", issues)
}
