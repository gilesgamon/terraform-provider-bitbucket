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

func dataIssueFieldMetrics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldMetricsRead,
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
			"metrics": {
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
						"value": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func dataIssueFieldMetricsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldMetricsRead", dumpResourceData(d, dataIssueFieldMetrics().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/metrics", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field metrics call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field metrics with params (%s): ", dumpResourceData(d, dataIssueFieldMetrics().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	metricsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field metrics response: %v", metricsBody)

	var metricsResponse IssueFieldMetricsResponse
	decodeerr := json.Unmarshal(metricsBody, &metricsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/metrics", workspace, repoSlug, fieldUUID))
	flattenIssueFieldMetrics(&metricsResponse, d)
	return nil
}

// IssueFieldMetricsResponse represents the response from the issue field metrics API
type IssueFieldMetricsResponse struct {
	Values []IssueFieldMetric `json:"values"`
	Page   int                `json:"page"`
	Size   int                `json:"size"`
	Next   string             `json:"next"`
}

// IssueFieldMetric represents a metric for an issue field
type IssueFieldMetric struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Value     float64                `json:"value"`
	Unit      string                 `json:"unit"`
	Timestamp string                 `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field metrics information
func flattenIssueFieldMetrics(c *IssueFieldMetricsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	metrics := make([]interface{}, len(c.Values))
	for i, metric := range c.Values {
		metrics[i] = map[string]interface{}{
			"uuid":      metric.UUID,
			"name":      metric.Name,
			"type":      metric.Type,
			"value":     metric.Value,
			"unit":      metric.Unit,
			"timestamp": metric.Timestamp,
			"metadata":  metric.Metadata,
			"links":     metric.Links,
		}
	}

	d.Set("metrics", metrics)
}
