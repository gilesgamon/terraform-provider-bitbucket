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

func dataPipelines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelinesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter pipelines by state (pending, in_progress, completed, error, stopped)",
			},
			"target": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter pipelines by target (commit, tag, branch, custom)",
			},
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Page number for pagination",
			},
			"pipelines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"build_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"trigger": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"target": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"completed_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duration_in_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"build_seconds_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"first_successful": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"expired": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"repository": {
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

func dataPipelinesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelinesRead", dumpResourceData(d, dataPipelines().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines", workspace, repoSlug)

	// Build query parameters
	params := make(map[string]string)
	if state, ok := d.GetOk("state"); ok {
		params["state"] = state.(string)
	}
	if target, ok := d.GetOk("target"); ok {
		params["target"] = target.(string)
	}
	if page, ok := d.GetOk("page"); ok {
		params["page"] = fmt.Sprintf("%d", page.(int))
	}

	// Add query parameters to URL
	if len(params) > 0 {
		url += "?"
		first := true
		for key, value := range params {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", key, value)
			first = false
		}
	}

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipelines call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipelines with params (%s): ", dumpResourceData(d, dataPipelines().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	pipelinesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipelines response: %v", pipelinesBody)

	var pipelinesResponse PipelinesResponse
	decodeerr := json.Unmarshal(pipelinesBody, &pipelinesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines", workspace, repoSlug))
	flattenPipelines(&pipelinesResponse, d)
	return nil
}

// PipelinesResponse represents the response from the pipelines API
type PipelinesResponse struct {
	Values []PipelineListItem `json:"values"`
	Page   int                `json:"page"`
	Size   int                `json:"size"`
	Next   string             `json:"next"`
}

// PipelineListItem represents a pipeline in the list
type PipelineListItem struct {
	UUID                string                 `json:"uuid"`
	BuildNumber         int                    `json:"build_number"`
	State               map[string]interface{} `json:"state"`
	Trigger             map[string]interface{} `json:"trigger"`
	Target              map[string]interface{} `json:"target"`
	CreatedOn           string                 `json:"created_on"`
	CompletedOn         string                 `json:"completed_on"`
	DurationInSeconds   int                    `json:"duration_in_seconds"`
	BuildSecondsUsed    int                    `json:"build_seconds_used"`
	FirstSuccessful     bool                   `json:"first_successful"`
	Expired             bool                   `json:"expired"`
	Repository          map[string]interface{} `json:"repository"`
}

// Flattens the pipelines information
func flattenPipelines(c *PipelinesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	pipelines := make([]interface{}, len(c.Values))
	for i, pipeline := range c.Values {
		pipelines[i] = map[string]interface{}{
			"uuid":                  pipeline.UUID,
			"build_number":          pipeline.BuildNumber,
			"state":                 pipeline.State,
			"trigger":               pipeline.Trigger,
			"target":                pipeline.Target,
			"created_on":            pipeline.CreatedOn,
			"completed_on":          pipeline.CompletedOn,
			"duration_in_seconds":   pipeline.DurationInSeconds,
			"build_seconds_used":    pipeline.BuildSecondsUsed,
			"first_successful":      pipeline.FirstSuccessful,
			"expired":               pipeline.Expired,
			"repository":            pipeline.Repository,
		}
	}

	d.Set("pipelines", pipelines)
}
