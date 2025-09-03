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

func dataBranchRestrictions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataBranchRestrictionsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"restrictions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"users": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
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
				},
			},
		},
	}
}

func dataBranchRestrictionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataBranchRestrictionsRead", dumpResourceData(d, dataBranchRestrictions().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/branch-restrictions", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from branch restrictions call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading branch restrictions with params (%s): ", dumpResourceData(d, dataBranchRestrictions().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	restrictionsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] branch restrictions response: %v", restrictionsBody)

	var restrictionsResponse BranchRestrictionsResponse
	decodeerr := json.Unmarshal(restrictionsBody, &restrictionsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/branch-restrictions", workspace, repoSlug))
	flattenBranchRestrictions(&restrictionsResponse, d)
	return nil
}

// BranchRestrictionsResponse represents the response from the branch restrictions API
type BranchRestrictionsResponse struct {
	Values []BranchRestrictionData `json:"values"`
	Page   int                     `json:"page"`
	Size   int                     `json:"size"`
	Next   string                  `json:"next"`
}

// BranchRestrictionData represents a branch restriction rule
type BranchRestrictionData struct {
	ID        int      `json:"id"`
	Kind      string   `json:"kind"`
	Pattern   string   `json:"pattern"`
	Value     int      `json:"value"`
	Enabled   bool     `json:"enabled"`
	Users     []string `json:"users"`
	Groups    []string `json:"groups"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the branch restrictions information
func flattenBranchRestrictions(c *BranchRestrictionsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	restrictions := make([]interface{}, len(c.Values))
	for i, restriction := range c.Values {
		restrictions[i] = map[string]interface{}{
			"id":       restriction.ID,
			"kind":     restriction.Kind,
			"pattern":  restriction.Pattern,
			"value":    restriction.Value,
			"enabled":  restriction.Enabled,
			"users":    restriction.Users,
			"groups":   restriction.Groups,
			"links":    restriction.Links,
		}
	}

	d.Set("restrictions", restrictions)
}
