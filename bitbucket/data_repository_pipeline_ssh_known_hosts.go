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

func dataRepositoryPipelineSSHKnownHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryPipelineSSHKnownHostsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"known_hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
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

func dataRepositoryPipelineSSHKnownHostsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryPipelineSSHKnownHostsRead", dumpResourceData(d, dataRepositoryPipelineSSHKnownHosts().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/ssh/known_hosts", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository pipeline SSH known hosts call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository pipeline SSH known hosts with params (%s): ", dumpResourceData(d, dataRepositoryPipelineSSHKnownHosts().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	knownHostsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository pipeline SSH known hosts response: %v", knownHostsBody)

	var knownHostsResponse RepositoryPipelineSSHKnownHostsResponse
	decodeerr := json.Unmarshal(knownHostsBody, &knownHostsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines_config/ssh/known_hosts", workspace, repoSlug))
	flattenRepositoryPipelineSSHKnownHosts(&knownHostsResponse, d)
	return nil
}

// RepositoryPipelineSSHKnownHostsResponse represents the response from the repository pipeline SSH known hosts API
type RepositoryPipelineSSHKnownHostsResponse struct {
	Values []RepositoryPipelineSSHKnownHost `json:"values"`
	Page   int                             `json:"page"`
	Size   int                             `json:"size"`
	Next   string                          `json:"next"`
}

// RepositoryPipelineSSHKnownHost represents an SSH known host
type RepositoryPipelineSSHKnownHost struct {
	UUID      string                 `json:"uuid"`
	Hostname  string                 `json:"hostname"`
	PublicKey string                 `json:"public_key"`
	Comment   string                 `json:"comment"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the repository pipeline SSH known hosts information
func flattenRepositoryPipelineSSHKnownHosts(c *RepositoryPipelineSSHKnownHostsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	knownHosts := make([]interface{}, len(c.Values))
	for i, host := range c.Values {
		knownHosts[i] = map[string]interface{}{
			"uuid":       host.UUID,
			"hostname":   host.Hostname,
			"public_key": host.PublicKey,
			"comment":    host.Comment,
			"created_on": host.CreatedOn,
			"updated_on": host.UpdatedOn,
			"links":      host.Links,
		}
	}

	d.Set("known_hosts", knownHosts)
}
