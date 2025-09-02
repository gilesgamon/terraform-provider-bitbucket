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
			"commit_sha": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Commit SHA to retrieve approvals for",
			},
			"approvals": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"approval_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"approver": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"approval_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"approved_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"branch_restriction": {
							Type:     schema.TypeString,
							Computed: true,
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
	commitSha := d.Get("commit_sha").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitApprovalsRead", dumpResourceData(d, dataCommitApprovals().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commit/%s/approvals",
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
		return diag.Errorf("no response returned from commit approvals call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commitSha, workspace, repoSlug)
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

	d.SetId(fmt.Sprintf("%s/%s/%s/approvals", workspace, repoSlug, commitSha))
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

// CommitApproval represents an approval for a commit
type CommitApproval struct {
	ApprovalID        string                 `json:"approval_id"`
	Approver          string                 `json:"approver"`
	ApprovalType      string                 `json:"approval_type"`
	Status            string                 `json:"status"`
	ApprovedOn        string                 `json:"approved_on"`
	Comment           string                 `json:"comment"`
	Required          bool                   `json:"required"`
	BranchRestriction string                 `json:"branch_restriction"`
	Links             map[string]interface{} `json:"links"`
}

// Flattens the commit approvals information
func flattenCommitApprovals(c *CommitApprovalsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	approvals := make([]interface{}, len(c.Values))
	for i, approval := range c.Values {
		approvals[i] = map[string]interface{}{
			"approval_id":        approval.ApprovalID,
			"approver":           approval.Approver,
			"approval_type":      approval.ApprovalType,
			"status":             approval.Status,
			"approved_on":        approval.ApprovedOn,
			"comment":            approval.Comment,
			"required":           approval.Required,
			"branch_restriction": approval.BranchRestriction,
		}
	}

	d.Set("approvals", approvals)
}
