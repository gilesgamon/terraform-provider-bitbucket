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

func dataIssueFieldValidations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueFieldValidationsRead,
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
			"validations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
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

func dataIssueFieldValidationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	fieldUUID := d.Get("field_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueFieldValidationsRead", dumpResourceData(d, dataIssueFieldValidations().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issue-fields/%s/validations", workspace, repoSlug, fieldUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue field validations call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue field %s in repository %s/%s", fieldUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue field validations with params (%s): ", dumpResourceData(d, dataIssueFieldValidations().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	validationsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue field validations response: %v", validationsBody)

	var validationsResponse IssueFieldValidationsResponse
	decodeerr := json.Unmarshal(validationsBody, &validationsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issue-fields/%s/validations", workspace, repoSlug, fieldUUID))
	flattenIssueFieldValidations(&validationsResponse, d)
	return nil
}

// IssueFieldValidationsResponse represents the response from the issue field validations API
type IssueFieldValidationsResponse struct {
	Values []IssueFieldValidation `json:"values"`
	Page   int                    `json:"page"`
	Size   int                    `json:"size"`
	Next   string                 `json:"next"`
}

// IssueFieldValidation represents a validation for an issue field
type IssueFieldValidation struct {
	UUID      string                 `json:"uuid"`
	Type      string                 `json:"type"`
	Value     string                 `json:"value"`
	Message   string                 `json:"message"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the issue field validations information
func flattenIssueFieldValidations(c *IssueFieldValidationsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	validations := make([]interface{}, len(c.Values))
	for i, validation := range c.Values {
		validations[i] = map[string]interface{}{
			"uuid":       validation.UUID,
			"type":       validation.Type,
			"value":      validation.Value,
			"message":    validation.Message,
			"created_on": validation.CreatedOn,
			"updated_on": validation.UpdatedOn,
			"links":      validation.Links,
		}
	}

	d.Set("validations", validations)
}
