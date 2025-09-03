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

func dataIssueAttachments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueAttachmentsRead,
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
				Type:     schema.TypeString,
				Required: true,
			},
			"attachments": {
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
						"size": {
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
				},
			},
		},
	}
}

func dataIssueAttachmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	issueID := d.Get("issue_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueAttachmentsRead", dumpResourceData(d, dataIssueAttachments().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues/%s/attachments", workspace, repoSlug, issueID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue attachments call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue %s in repository %s/%s", issueID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue attachments with params (%s): ", dumpResourceData(d, dataIssueAttachments().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	attachmentsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue attachments response: %v", attachmentsBody)

	var attachmentsResponse IssueAttachmentsResponse
	decodeerr := json.Unmarshal(attachmentsBody, &attachmentsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues/%s/attachments", workspace, repoSlug, issueID))
	flattenIssueAttachments(&attachmentsResponse, d)
	return nil
}

// IssueAttachmentsResponse represents the response from the issue attachments API
type IssueAttachmentsResponse struct {
	Values []IssueAttachment `json:"values"`
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Next   string            `json:"next"`
}

// IssueAttachment represents an attachment on an issue
type IssueAttachment struct {
	UUID string                 `json:"uuid"`
	Name string                 `json:"name"`
	Size int                    `json:"size"`
	Links map[string]interface{} `json:"links"`
}

// Flattens the issue attachments information
func flattenIssueAttachments(c *IssueAttachmentsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	attachments := make([]interface{}, len(c.Values))
	for i, attachment := range c.Values {
		attachments[i] = map[string]interface{}{
			"uuid":  attachment.UUID,
			"name":  attachment.Name,
			"size":  attachment.Size,
			"links": attachment.Links,
		}
	}

	d.Set("attachments", attachments)
}
