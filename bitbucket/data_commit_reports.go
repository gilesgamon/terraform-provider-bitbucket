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

func dataCommitReports() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCommitReportsRead,
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
				Description: "Commit SHA to retrieve reports for",
			},
			"reports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"report_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"details": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"report_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reporter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
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
						"annotations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"annotation_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"line": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"severity": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataCommitReportsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	commitSha := d.Get("commit_sha").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataCommitReportsRead", dumpResourceData(d, dataCommitReports().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/commit/%s/reports",
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
		return diag.Errorf("no response returned from commit reports call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate commit %s in repository %s/%s", commitSha, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading commit reports with params (%s): ", dumpResourceData(d, dataCommitReports().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	reportsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] commit reports response: %v", reportsBody)

	var reportsResponse CommitReportsResponse
	decodeerr := json.Unmarshal(reportsBody, &reportsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/reports", workspace, repoSlug, commitSha))
	flattenCommitReports(&reportsResponse, d)
	return nil
}

// CommitReportsResponse represents the response from the commit reports API
type CommitReportsResponse struct {
	Values []CommitReport `json:"values"`
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Next   string         `json:"next"`
}

// CommitReport represents a report for a commit
type CommitReport struct {
	ReportID   string                 `json:"report_id"`
	Title      string                 `json:"title"`
	Details    string                 `json:"details"`
	ReportType string                 `json:"report_type"`
	Reporter   string                 `json:"reporter"`
	Result     string                 `json:"result"`
	Severity   string                 `json:"severity"`
	CreatedOn  string                 `json:"created_on"`
	UpdatedOn  string                 `json:"updated_on"`
	Links      map[string]interface{} `json:"links"`
}

// Flattens the commit reports information
func flattenCommitReports(c *CommitReportsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	reports := make([]interface{}, len(c.Values))
	for i, report := range c.Values {
		reports[i] = map[string]interface{}{
			"report_id":   report.ReportID,
			"title":       report.Title,
			"details":     report.Details,
			"report_type": report.ReportType,
			"reporter":    report.Reporter,
			"result":      report.Result,
			"severity":    report.Severity,
			"created_on":  report.CreatedOn,
			"updated_on":  report.UpdatedOn,
		}
	}

	d.Set("reports", reports)
}
