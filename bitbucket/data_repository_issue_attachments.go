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

func dataRepositoryIssueAttachments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryIssueAttachmentsRead,
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Issue ID",
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_on": {
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

func dataRepositoryIssueAttachmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	issueID := d.Get("issue_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryIssueAttachmentsRead", dumpResourceData(d, dataRepositoryIssueAttachments().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues/%s/attachments", workspace, repoSlug, issueID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository issue attachments call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue %s in repository %s/%s", issueID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository issue attachments with params (%s): ", dumpResourceData(d, dataRepositoryIssueAttachments().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	attachmentsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository issue attachments response: %v", attachmentsBody)

	var attachmentsResponse RepositoryIssueAttachmentsResponse
	decodeerr := json.Unmarshal(attachmentsBody, &attachmentsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues/%s/attachments", workspace, repoSlug, issueID))
	flattenRepositoryIssueAttachments(&attachmentsResponse, d)
	return nil
}

// RepositoryIssueAttachmentsResponse represents the response from the repository issue attachments API
type RepositoryIssueAttachmentsResponse struct {
	Values []RepositoryIssueAttachment `json:"values"`
	Page   int                         `json:"page"`
	Size   int                         `json:"size"`
	Next   string                      `json:"next"`
}

// RepositoryIssueAttachment represents an issue attachment
type RepositoryIssueAttachment struct {
	Name      string                 `json:"name"`
	Size      int                    `json:"size"`
	CreatedOn string                 `json:"created_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the repository issue attachments information
func flattenRepositoryIssueAttachments(c *RepositoryIssueAttachmentsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	attachments := make([]interface{}, len(c.Values))
	for i, attachment := range c.Values {
		attachments[i] = map[string]interface{}{
			"name":       attachment.Name,
			"size":       attachment.Size,
			"created_on": attachment.CreatedOn,
			"links":      attachment.Links,
		}
	}

	d.Set("attachments", attachments)
}
