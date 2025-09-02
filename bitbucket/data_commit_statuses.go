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

func dataCommitStatuses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCommitStatusesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commit_sha": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Commit SHA to retrieve statuses for",
			},
			"statuses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"refname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commit": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataCommitStatusesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	commitSha := d.Get("commit_sha").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitStatusesRead", dumpResourceData(d, dataCommitStatuses().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commit/%s/statuses",
		workspace,
		repoSlug,
		commitSha,
	)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from commit statuses call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commitSha, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading commit statuses with params (%s): ", dumpResourceData(d, dataCommitStatuses().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	statusesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] commit statuses response: %v", statusesBody)

	var statusesResponse CommitStatusesResponse
	decodeerr := json.Unmarshal(statusesBody, &statusesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/statuses", workspace, repoSlug, commitSha))
	flattenCommitStatuses(&statusesResponse, d)
	return nil
}

// CommitStatusesResponse represents the response from the commit statuses API
type CommitStatusesResponse struct {
	Values []CommitStatus `json:"values"`
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Next   string         `json:"next"`
}

// CommitStatus represents a status for a commit (build, test, etc.)
type CommitStatus struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	URL         string                 `json:"url"`
	State       string                 `json:"state"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Type        string                 `json:"type"`
	Refname     string                 `json:"refname"`
	Commit      string                 `json:"commit"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the commit statuses information
func flattenCommitStatuses(c *CommitStatusesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	statuses := make([]interface{}, len(c.Values))
	for i, status := range c.Values {
		statuses[i] = map[string]interface{}{
			"key":         status.Key,
			"name":        status.Name,
			"description": status.Description,
			"url":         status.URL,
			"state":       status.State,
			"created_on":  status.CreatedOn,
			"updated_on":  status.UpdatedOn,
			"type":        status.Type,
			"refname":     status.Refname,
			"commit":      status.Commit,
		}
	}

	d.Set("statuses", statuses)
}
