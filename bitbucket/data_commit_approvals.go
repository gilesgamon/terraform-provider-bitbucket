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

func dataCommitApprovals() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCommitApprovalsRead,
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
			"approvals": {
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
						"approved": {
							Type:     schema.TypeBool,
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

func dataCommitApprovalsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	commit := d.Get("commit").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitApprovalsRead", dumpResourceData(d, dataCommitApprovals().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commits/%s/approvals", workspace, repoSlug, commit)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from commit approvals call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commit, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading commit approvals with params (%s): ", dumpResourceData(d, dataCommitApprovals().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	approvalsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] commit approvals response: %v", approvalsBody)

	var approvalsResponse CommitApprovalsResponse
	decodeerr := json.Unmarshal(approvalsBody, &approvalsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/commits/%s/approvals", workspace, repoSlug, commit))
	flattenCommitApprovals(&approvalsResponse, d)
	return nil
}

// CommitApprovalsResponse represents the response from the commit approvals API
type CommitApprovalsResponse struct {
	Values []CommitApproval `json:"values"`
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Next   string           `json:"next"`
}

// CommitApproval represents a commit approval
type CommitApproval struct {
	UUID      string                 `json:"uuid"`
	User      map[string]interface{} `json:"user"`
	Approved  bool                   `json:"approved"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the commit approvals information
func flattenCommitApprovals(c *CommitApprovalsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	approvals := make([]interface{}, len(c.Values))
	for i, approval := range c.Values {
		approvals[i] = map[string]interface{}{
			"uuid":       approval.UUID,
			"user":       approval.User,
			"approved":   approval.Approved,
			"created_on": approval.CreatedOn,
			"updated_on": approval.UpdatedOn,
			"links":      approval.Links,
		}
	}

	d.Set("approvals", approvals)
}
