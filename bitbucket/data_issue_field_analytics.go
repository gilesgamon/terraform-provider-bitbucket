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

func dataIssueFieldAnalytics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldAnalyticsRead,
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
			"analytics": {
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
						"data": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"insights": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"trends": {
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

func dataIssueFieldAnalyticsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldAnalyticsRead", dumpResourceData(d, dataIssueFieldAnalytics().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/analytics", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field analytics call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field analytics with params (%s): ", dumpResourceData(d, dataIssueFieldAnalytics().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	analyticsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field analytics response: %v", analyticsBody)

	var analyticsResponse IssueFieldAnalyticsResponse
	decodeerr := json.Unmarshal(analyticsBody, &analyticsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/analytics", workspace, repoSlug, fieldUUID))
	flattenIssueFieldAnalytics(&analyticsResponse, d)
	return nil
}

// IssueFieldAnalyticsResponse represents the response from the issue field analytics API
type IssueFieldAnalyticsResponse struct {
	Values []IssueFieldAnalytic `json:"values"`
	Page   int                  `json:"page"`
	Size   int                  `json:"size"`
	Next   string               `json:"next"`
}

// IssueFieldAnalytic represents analytics for an issue field
type IssueFieldAnalytic struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Insights  []string               `json:"insights"`
	Trends    map[string]interface{} `json:"trends"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field analytics information
func flattenIssueFieldAnalytics(c *IssueFieldAnalyticsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	analytics := make([]interface{}, len(c.Values))
	for i, analytic := range c.Values {
		analytics[i] = map[string]interface{}{
			"uuid":       analytic.UUID,
			"name":       analytic.Name,
			"type":       analytic.Type,
			"data":       analytic.Data,
			"insights":   analytic.Insights,
			"trends":     analytic.Trends,
			"created_on": analytic.CreatedOn,
			"updated_on": analytic.UpdatedOn,
			"links":      analytic.Links,
		}
	}

	d.Set("analytics", analytics)
}
