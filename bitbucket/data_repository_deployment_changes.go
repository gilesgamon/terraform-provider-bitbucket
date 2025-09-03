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

func dataRepositoryDeploymentChanges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryDeploymentChangesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"changes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deployment": {
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

func dataRepositoryDeploymentChangesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	environmentUUID := d.Get("environment_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryDeploymentChangesRead", dumpResourceData(d, dataRepositoryDeploymentChanges().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/environments/%s/changes", workspace, repoSlug, environmentUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository deployment changes call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate environment %s for repository %s/%s", environmentUUID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository deployment changes with params (%s): ", dumpResourceData(d, dataRepositoryDeploymentChanges().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	changesBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository deployment changes response: %v", changesBody)

	var changesResponse RepositoryDeploymentChangesResponse
	decodeerr := json.Unmarshal(changesBody, &changesResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/environments/%s/changes", workspace, repoSlug, environmentUUID))
	flattenRepositoryDeploymentChanges(&changesResponse, d)
	return nil
}

// RepositoryDeploymentChangesResponse represents the response from the repository deployment changes API
type RepositoryDeploymentChangesResponse struct {
	Values []RepositoryDeploymentChange `json:"values"`
	Page   int                          `json:"page"`
	Size   int                          `json:"size"`
	Next   string                       `json:"next"`
}

// RepositoryDeploymentChange represents a deployment change
type RepositoryDeploymentChange struct {
	UUID        string                 `json:"uuid"`
	Version     string                 `json:"version"`
	State       string                 `json:"state"`
	Deployment  map[string]interface{} `json:"deployment"`
	CreatedOn   string                 `json:"created_on"`
	UpdatedOn   string                 `json:"updated_on"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the repository deployment changes information
func flattenRepositoryDeploymentChanges(c *RepositoryDeploymentChangesResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	changes := make([]interface{}, len(c.Values))
	for i, change := range c.Values {
		changes[i] = map[string]interface{}{
			"uuid":       change.UUID,
			"version":    change.Version,
			"state":      change.State,
			"deployment": change.Deployment,
			"created_on": change.CreatedOn,
			"updated_on": change.UpdatedOn,
			"links":      change.Links,
		}
	}

	d.Set("changes", changes)
}
