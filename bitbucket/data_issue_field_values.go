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

func dataIssueFieldValues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldValuesRead,
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
			"values": {
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

func dataIssueFieldValuesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldValuesRead", dumpResourceData(d, dataIssueFieldValues().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/values", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field values call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field values with params (%s): ", dumpResourceData(d, dataIssueFieldValues().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	valuesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field values response: %v", valuesBody)

	var valuesResponse IssueFieldValuesResponse
	decodeerr := json.Unmarshal(valuesBody, &valuesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/values", workspace, repoSlug, fieldUUID))
	flattenIssueFieldValues(&valuesResponse, d)
	return nil
}

// IssueFieldValuesResponse represents the response from the issue field values API
type IssueFieldValuesResponse struct {
	Values []IssueFieldValue `json:"values"`
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Next   string            `json:"next"`
}

// IssueFieldValue represents a value for an issue field
type IssueFieldValue struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Value     string                 `json:"value"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field values information
func flattenIssueFieldValues(c *IssueFieldValuesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	values := make([]interface{}, len(c.Values))
	for i, value := range c.Values {
		values[i] = map[string]interface{}{
			"uuid":       value.UUID,
			"name":       value.Name,
			"value":      value.Value,
			"created_on": value.CreatedOn,
			"updated_on": value.UpdatedOn,
			"links":      value.Links,
		}
	}

	d.Set("values", values)
}
