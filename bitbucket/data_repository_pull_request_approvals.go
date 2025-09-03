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

func dataRepositoryPullRequestApprovals() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryPullRequestApprovalsRead,
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pull request ID",
			},
			"approvals": {
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
						"approved_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
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

func dataRepositoryPullRequestApprovalsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pullRequestID := d.Get("pull_request_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryPullRequestApprovalsRead", dumpResourceData(d, dataRepositoryPullRequestApprovals().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pullrequests/%s/approve", workspace, repoSlug, pullRequestID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository pull request approvals call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pull request %s in repository %s/%s", pullRequestID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository pull request approvals with params (%s): ", dumpResourceData(d, dataRepositoryPullRequestApprovals().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	approvalsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository pull request approvals response: %v", approvalsBody)

	var approvalsResponse RepositoryPullRequestApprovalsResponse
	decodeerr := json.Unmarshal(approvalsBody, &approvalsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pullrequests/%s/approve", workspace, repoSlug, pullRequestID))
	flattenRepositoryPullRequestApprovals(&approvalsResponse, d)
	return nil
}

// RepositoryPullRequestApprovalsResponse represents the response from the repository pull request approvals API
type RepositoryPullRequestApprovalsResponse struct {
	Values []RepositoryPullRequestApproval `json:"values"`
	Page   int                             `json:"page"`
	Size   int                             `json:"size"`
	Next   string                          `json:"next"`
}

// RepositoryPullRequestApproval represents a pull request approval
type RepositoryPullRequestApproval struct {
	User        map[string]interface{} `json:"user"`
	ApprovedOn  string                 `json:"approved_on"`
	Role        string                 `json:"role"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the repository pull request approvals information
func flattenRepositoryPullRequestApprovals(c *RepositoryPullRequestApprovalsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	approvals := make([]interface{}, len(c.Values))
	for i, approval := range c.Values {
		approvals[i] = map[string]interface{}{
			"user":         approval.User,
			"approved_on":  approval.ApprovedOn,
			"role":         approval.Role,
			"links":        approval.Links,
		}
	}

	d.Set("approvals", approvals)
}
