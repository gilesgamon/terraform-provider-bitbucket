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

func dataPipelineVariables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineVariablesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline_uuid": {
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
						"type": {
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

func dataPipelineVariablesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pipelineUUID := d.Get("pipeline_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineVariablesRead", dumpResourceData(d, dataPipelineVariables().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/%s/variables", workspace, repoSlug, pipelineUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline variables call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate pipeline %s in repository %s/%s", pipelineUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline variables with params (%s): ", dumpResourceData(d, dataPipelineVariables().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	variablesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline variables response: %v", variablesBody)

	var variablesResponse PipelineVariablesResponse
	decodeerr := json.Unmarshal(variablesBody, &variablesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/%s/variables", workspace, repoSlug, pipelineUUID))
	flattenPipelineVariables(&variablesResponse, d)
	return nil
}

// PipelineVariablesResponse represents the response from the pipeline variables API
type PipelineVariablesResponse struct {
	Values []PipelineVariable `json:"values"`
	Page   int                `json:"page"`
	Size   int                `json:"size"`
	Next   string             `json:"next"`
}

// PipelineVariable represents a pipeline variable
type PipelineVariable struct {
	UUID      string                 `json:"uuid"`
	Key       string                 `json:"key"`
	Value     string                 `json:"value"`
	Secured   bool                   `json:"secured"`
	Type      string                 `json:"type"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the pipeline variables information
func flattenPipelineVariables(c *PipelineVariablesResponse, d *schema.ResourceData) {
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
			"type":       variable.Type,
			"created_on": variable.CreatedOn,
			"updated_on": variable.UpdatedOn,
			"links":      variable.Links,
		}
	}

	d.Set("variables", variables)
}
