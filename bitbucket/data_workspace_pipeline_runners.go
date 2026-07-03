package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWorkspacePipelineRunners() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataWorkspacePipelineRunnersRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"runners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: pipelineRunnerSchema(),
				},
			},
		},
	}
}

func dataWorkspacePipelineRunnersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataWorkspacePipelineRunnersRead", dumpResourceData(d, dataWorkspacePipelineRunners().Schema))

	url := fmt.Sprintf("2.0/workspaces/%s/pipelines-config/runners", workspace)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from workspace pipeline runners call. Make sure your credentials are accurate.")
	}

	if res.Body == nil {
		return diag.Errorf("error reading workspace pipeline runners with params (%s): ", dumpResourceData(d, dataWorkspacePipelineRunners().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	runnersBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)

	var runnersResponse PipelineRunnersResponse
	decodeerr := json.Unmarshal(runnersBody, &runnersResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/pipelines-config/runners", workspace))
	d.Set("runners", flattenPipelineRunners(runnersResponse.Values))
	return nil
}

// PipelineRunnersResponse represents a paginated list of pipeline runners
type PipelineRunnersResponse struct {
	Values  []PipelineRunner `json:"values"`
	Page    int              `json:"page"`
	Size    int              `json:"size"`
	Pagelen int              `json:"pagelen"`
	Next    string           `json:"next"`
}

// PipelineRunner represents a Bitbucket Pipelines self-hosted runner
type PipelineRunner struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Labels      []string               `json:"labels"`
	State       *PipelineRunnerState   `json:"state"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	OauthClient map[string]interface{} `json:"oauth_client"`
}

// PipelineRunnerState represents the state information of a runner
type PipelineRunnerState struct {
	Status    string                 `json:"status"`
	Version   map[string]interface{} `json:"version"`
	UpdatedOn string                 `json:"updated_on"`
	Cordoned  bool                   `json:"cordoned"`
}

// pipelineRunnerSchema returns the shared schema describing a single runner
func pipelineRunnerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"labels": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"created_on": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_on": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"cordoned": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"updated_on": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"version": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"oauth_client": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

// flattenPipelineRunners converts a list of runners into a Terraform-friendly structure
func flattenPipelineRunners(runners []PipelineRunner) []interface{} {
	result := make([]interface{}, len(runners))
	for i, runner := range runners {
		result[i] = flattenPipelineRunner(&runner)
	}
	return result
}

// flattenPipelineRunner converts a single runner into a Terraform-friendly map
func flattenPipelineRunner(runner *PipelineRunner) map[string]interface{} {
	if runner == nil {
		return nil
	}

	labels := make([]interface{}, len(runner.Labels))
	for i, l := range runner.Labels {
		labels[i] = l
	}

	m := map[string]interface{}{
		"uuid":         runner.UUID,
		"name":         runner.Name,
		"labels":       labels,
		"created_on":   runner.CreatedOn,
		"updated_on":   runner.UpdatedOn,
		"oauth_client": stringifyMap(runner.OauthClient),
	}

	if runner.State != nil {
		m["state"] = []interface{}{
			map[string]interface{}{
				"status":     runner.State.Status,
				"cordoned":   runner.State.Cordoned,
				"updated_on": runner.State.UpdatedOn,
				"version":    stringifyMap(runner.State.Version),
			},
		}
	} else {
		m["state"] = []interface{}{}
	}

	return m
}

// stringifyMap converts an arbitrary map into a map[string]string so it can be
// stored in a Terraform TypeMap attribute.
func stringifyMap(in map[string]interface{}) map[string]interface{} {
	if in == nil {
		return map[string]interface{}{}
	}
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		switch val := v.(type) {
		case string:
			out[k] = val
		default:
			b, err := json.Marshal(val)
			if err != nil {
				out[k] = fmt.Sprintf("%v", val)
			} else {
				out[k] = string(b)
			}
		}
	}
	return out
}
