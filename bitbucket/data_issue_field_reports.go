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

func dataIssueFieldReports() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldReportsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"field_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data": {
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

func dataIssueFieldReportsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldReportsRead", dumpResourceData(d, dataIssueFieldReports().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/reports", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field reports call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field reports with params (%s): ", dumpResourceData(d, dataIssueFieldReports().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	reportsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field reports response: %v", reportsBody)

	var reportsResponse IssueFieldReportsResponse
	decodeerr := json.Unmarshal(reportsBody, &reportsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/reports", workspace, repoSlug, fieldUUID))
	flattenIssueFieldReports(&reportsResponse, d)
	return nil
}

// IssueFieldReportsResponse represents the response from the issue field reports API
type IssueFieldReportsResponse struct {
	Values []IssueFieldReport `json:"values"`
	Page   int                `json:"page"`
	Size   int                `json:"size"`
	Next   string             `json:"next"`
}

// IssueFieldReport represents a report for an issue field
type IssueFieldReport struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Format    string                 `json:"format"`
	Data      map[string]interface{} `json:"data"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field reports information
func flattenIssueFieldReports(c *IssueFieldReportsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	reports := make([]interface{}, len(c.Values))
	for i, report := range c.Values {
		reports[i] = map[string]interface{}{
			"uuid":       report.UUID,
			"name":       report.Name,
			"type":       report.Type,
			"format":     report.Format,
			"data":       report.Data,
			"created_on": report.CreatedOn,
			"updated_on": report.UpdatedOn,
			"links":      report.Links,
		}
	}

	d.Set("reports", reports)
}
