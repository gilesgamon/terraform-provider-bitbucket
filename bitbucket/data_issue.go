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

func dataIssue() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueRead,
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
				Description: "Issue ID to retrieve",
			},
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
	}
}

func dataIssueRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	issueID := d.Get("issue_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueRead", dumpResourceData(d, dataIssue().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues/%s", workspace, repoSlug, issueID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue %s in repository %s/%s", issueID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue with params (%s): ", dumpResourceData(d, dataIssue().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	issueBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue response: %v", issueBody)

	var issue IssueDetail
	decodeerr := json.Unmarshal(issueBody, &issue)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues/%s", workspace, repoSlug, issueID))
	flattenIssue(&issue, d)
	return nil
}

// IssueDetail represents a single issue from the API
type IssueDetail struct {
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

// Flattens the issue information
func flattenIssue(i *IssueDetail, d *schema.ResourceData) {
	if i == nil {
		return
	}

	d.Set("id", i.ID)
	d.Set("title", i.Title)
	d.Set("content", i.Content)
	d.Set("state", i.State)
	d.Set("kind", i.Kind)
	d.Set("priority", i.Priority)
	d.Set("assignee", i.Assignee)
	d.Set("reporter", i.Reporter)
	d.Set("milestone", i.Milestone)
	d.Set("component", i.Component)
	d.Set("version", i.Version)
	d.Set("created_on", i.CreatedOn)
	d.Set("updated_on", i.UpdatedOn)
	d.Set("votes", i.Votes)
	d.Set("watches", i.Watches)
	d.Set("links", i.Links)
}
