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

func dataRepositoryVariables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryVariablesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"variables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secured": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"uuid": {
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
					},
				},
			},
		},
	}
}

func dataRepositoryVariablesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryVariablesRead", dumpResourceData(d, dataRepositoryVariables().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/variables", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository variables call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository variables with params (%s): ", dumpResourceData(d, dataRepositoryVariables().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	variablesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository variables response: %v", variablesBody)

	var variablesResponse RepositoryVariablesResponse
	decodeerr := json.Unmarshal(variablesBody, &variablesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/variables", workspace, repoSlug))
	flattenRepositoryVariables(&variablesResponse, d)
	return nil
}

// RepositoryVariablesResponse represents the response from the repository variables API
type RepositoryVariablesResponse struct {
	Values []RepositoryVariable `json:"values"`
	Page   int                  `json:"page"`
	Size   int                  `json:"size"`
	Next   string               `json:"next"`
}

// RepositoryVariable represents a variable in a repository
type RepositoryVariable struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Secured   bool   `json:"secured"`
	UUID      string `json:"uuid"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}

// Flattens the repository variables information
func flattenRepositoryVariables(c *RepositoryVariablesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	variables := make([]interface{}, len(c.Values))
	for i, variable := range c.Values {
		variables[i] = map[string]interface{}{
			"key":        variable.Key,
			"value":      variable.Value,
			"secured":    variable.Secured,
			"uuid":       variable.UUID,
			"created_on": variable.CreatedOn,
			"updated_on": variable.UpdatedOn,
		}
	}

	d.Set("variables", variables)
}
