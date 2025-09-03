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

func dataEffectiveBranchingModel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataEffectiveBranchingModelRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"development": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"production": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"branch_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"links": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataEffectiveBranchingModelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataEffectiveBranchingModelRead", dumpResourceData(d, dataEffectiveBranchingModel().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/branching-model/effective", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from effective branching model call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading effective branching model with params (%s): ", dumpResourceData(d, dataEffectiveBranchingModel().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	modelBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] effective branching model response: %v", modelBody)

	var effectiveModel EffectiveBranchingModelData
	decodeerr := json.Unmarshal(modelBody, &effectiveModel)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/branching-model/effective", workspace, repoSlug))
	flattenEffectiveBranchingModel(&effectiveModel, d)
	return nil
}

// EffectiveBranchingModelData represents the effective branching model configuration
type EffectiveBranchingModelData struct {
	Development  map[string]interface{} `json:"development"`
	Production   map[string]interface{} `json:"production"`
	BranchTypes  []EffectiveBranchTypeData `json:"branch_types"`
	Links        map[string]interface{} `json:"links"`
}

// EffectiveBranchTypeData represents an effective branch type configuration
type EffectiveBranchTypeData struct {
	Kind   string `json:"kind"`
	Prefix string `json:"prefix"`
	Enabled bool  `json:"enabled"`
}

// Flattens the effective branching model information
func flattenEffectiveBranchingModel(c *EffectiveBranchingModelData, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("development", c.Development)
	d.Set("production", c.Production)
	d.Set("links", c.Links)

	branchTypes := make([]interface{}, len(c.BranchTypes))
	for i, branchType := range c.BranchTypes {
		branchTypes[i] = map[string]interface{}{
			"kind":   branchType.Kind,
			"prefix": branchType.Prefix,
			"enabled": branchType.Enabled,
		}
	}

	d.Set("branch_types", branchTypes)
}
