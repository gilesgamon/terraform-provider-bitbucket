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

func dataIssueFieldOptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldOptionsRead,
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
			"options": {
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
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"icon": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"color": {
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

func dataIssueFieldOptionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldOptionsRead", dumpResourceData(d, dataIssueFieldOptions().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/options", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field options call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field options with params (%s): ", dumpResourceData(d, dataIssueFieldOptions().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	optionsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field options response: %v", optionsBody)

	var optionsResponse IssueFieldOptionsResponse
	decodeerr := json.Unmarshal(optionsBody, &optionsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/options", workspace, repoSlug, fieldUUID))
	flattenIssueFieldOptions(&optionsResponse, d)
	return nil
}

// IssueFieldOptionsResponse represents the response from the issue field options API
type IssueFieldOptionsResponse struct {
	Values []IssueFieldOption `json:"values"`
	Page   int                `json:"page"`
	Size   int                `json:"size"`
	Next   string             `json:"next"`
}

// IssueFieldOption represents an option for an issue field
type IssueFieldOption struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Value     string                 `json:"value"`
	Icon      string                 `json:"icon"`
	Color     string                 `json:"color"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field options information
func flattenIssueFieldOptions(c *IssueFieldOptionsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	options := make([]interface{}, len(c.Values))
	for i, option := range c.Values {
		options[i] = map[string]interface{}{
			"uuid":       option.UUID,
			"name":       option.Name,
			"value":      option.Value,
			"icon":       option.Icon,
			"color":      option.Color,
			"created_on": option.CreatedOn,
			"updated_on": option.UpdatedOn,
			"links":      option.Links,
		}
	}

	d.Set("options", options)
}
