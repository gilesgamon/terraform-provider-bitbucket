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

func dataIssueFieldLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldLogsRead,
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
			"logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timestamp": {
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
						"details": {
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

func dataIssueFieldLogsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldLogsRead", dumpResourceData(d, dataIssueFieldLogs().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/logs", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field logs call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field logs with params (%s): ", dumpResourceData(d, dataIssueFieldLogs().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	logsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field logs response: %v", logsBody)

	var logsResponse IssueFieldLogsResponse
	decodeerr := json.Unmarshal(logsBody, &logsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/logs", workspace, repoSlug, fieldUUID))
	flattenIssueFieldLogs(&logsResponse, d)
	return nil
}

// IssueFieldLogsResponse represents the response from the issue field logs API
type IssueFieldLogsResponse struct {
	Values []IssueFieldLog `json:"values"`
	Page   int             `json:"page"`
	Size   int             `json:"size"`
	Next   string          `json:"next"`
}

// IssueFieldLog represents a log entry for an issue field
type IssueFieldLog struct {
	UUID      string                 `json:"uuid"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	User      map[string]interface{} `json:"user"`
	Details   map[string]interface{} `json:"details"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field logs information
func flattenIssueFieldLogs(c *IssueFieldLogsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	logs := make([]interface{}, len(c.Values))
	for i, logEntry := range c.Values {
		logs[i] = map[string]interface{}{
			"uuid":      logEntry.UUID,
			"level":     logEntry.Level,
			"message":   logEntry.Message,
			"timestamp": logEntry.Timestamp,
			"user":      logEntry.User,
			"details":   logEntry.Details,
			"links":     logEntry.Links,
		}
	}

	d.Set("logs", logs)
}
