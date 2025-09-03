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

func dataWorkspaceVariables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataWorkspaceVariablesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"variables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
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

func dataWorkspaceVariablesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataWorkspaceVariablesRead", dumpResourceData(d, dataWorkspaceVariables().Schema))

	url := fmt.Sprintf("2.0/workspaces/%s/pipelines-config/variables", workspace)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from workspace variables call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate workspace %s", workspace)
	}

	if res.Body == nil {
		return diag.Errorf("error reading workspace variables with params (%s): ", dumpResourceData(d, dataWorkspaceVariables().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	variablesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] workspace variables response: %v", variablesBody)

	var variablesResponse WorkspaceVariablesResponse
	decodeerr := json.Unmarshal(variablesBody, &variablesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/variables", workspace))
	flattenWorkspaceVariables(&variablesResponse, d)
	return nil
}

// WorkspaceVariablesResponse represents the response from the workspace variables API
type WorkspaceVariablesResponse struct {
	Values []WorkspaceVariable `json:"values"`
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	Next   string              `json:"next"`
}

// WorkspaceVariable represents a variable in a workspace
type WorkspaceVariable struct {
	UUID      string `json:"uuid"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Secured   bool   `json:"secured"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}

// Flattens the workspace variables information
func flattenWorkspaceVariables(c *WorkspaceVariablesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	variables := make([]interface{}, len(c.Values))
	for i, variable := range c.Values {
		variables[i] = map[string]interface{}{
			"uuid":       variable.UUID,
			"key":        variable.Key,
			"value":      variable.Value,
			"secured":    variable.Secured,
			"created_on": variable.CreatedOn,
			"updated_on": variable.UpdatedOn,
		}
	}

	d.Set("variables", variables)
}
