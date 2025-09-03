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

func dataIssueFields() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fields": {
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
						"required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"default_value": {
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

func dataIssueFieldsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldsRead", dumpResourceData(d, dataIssueFields().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue fields call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue fields with params (%s): ", dumpResourceData(d, dataIssueFields().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	fieldsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue fields response: %v", fieldsBody)

	var fieldsResponse IssueFieldsResponse
	decodeerr := json.Unmarshal(fieldsBody, &fieldsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields", workspace, repoSlug))
	flattenIssueFields(&fieldsResponse, d)
	return nil
}

// IssueFieldsResponse represents the response from the issue fields API
type IssueFieldsResponse struct {
	Values []IssueField `json:"values"`
	Page   int          `json:"page"`
	Size   int          `json:"size"`
	Next   string       `json:"next"`
}

// IssueField represents an issue field
type IssueField struct {
	UUID         string                 `json:"uuid"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Required     bool                   `json:"required"`
	DefaultValue string                 `json:"default_value"`
	CreatedOn    string                 `json:"created_on"`
	UpdatedOn    string                 `json:"updated_on"`
	Links        map[string]interface{} `json:"links"`
}

// Flattens the issue fields information
func flattenIssueFields(c *IssueFieldsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	fields := make([]interface{}, len(c.Values))
	for i, field := range c.Values {
		fields[i] = map[string]interface{}{
			"uuid":          field.UUID,
			"name":          field.Name,
			"type":          field.Type,
			"required":      field.Required,
			"default_value": field.DefaultValue,
			"created_on":    field.CreatedOn,
			"updated_on":    field.UpdatedOn,
			"links":         field.Links,
		}
	}

	d.Set("fields", fields)
}
