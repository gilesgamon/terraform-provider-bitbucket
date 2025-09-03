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

func dataRepositoryPipelineSSHKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryPipelineSSHKeysRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ssh_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"label": {
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

func dataRepositoryPipelineSSHKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryPipelineSSHKeysRead", dumpResourceData(d, dataRepositoryPipelineSSHKeys().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/ssh/keys", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository pipeline SSH keys call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository pipeline SSH keys with params (%s): ", dumpResourceData(d, dataRepositoryPipelineSSHKeys().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	sshKeysBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository pipeline SSH keys response: %v", sshKeysBody)

	var sshKeysResponse RepositoryPipelineSSHKeysResponse
	decodeerr := json.Unmarshal(sshKeysBody, &sshKeysResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines_config/ssh/keys", workspace, repoSlug))
	flattenRepositoryPipelineSSHKeys(&sshKeysResponse, d)
	return nil
}

// RepositoryPipelineSSHKeysResponse represents the response from the repository pipeline SSH keys API
type RepositoryPipelineSSHKeysResponse struct {
	Values []RepositoryPipelineSSHKey `json:"values"`
	Page   int                        `json:"page"`
	Size   int                        `json:"size"`
	Next   string                     `json:"next"`
}

// RepositoryPipelineSSHKey represents an SSH key
type RepositoryPipelineSSHKey struct {
	UUID      string                 `json:"uuid"`
	Label     string                 `json:"label"`
	PublicKey string                 `json:"public_key"`
	Comment   string                 `json:"comment"`
	CreatedOn string                 `json:"created_on"`
	UpdatedOn string                 `json:"updated_on"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the repository pipeline SSH keys information
func flattenRepositoryPipelineSSHKeys(c *RepositoryPipelineSSHKeysResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	sshKeys := make([]interface{}, len(c.Values))
	for i, sshKey := range c.Values {
		sshKeys[i] = map[string]interface{}{
			"uuid":       sshKey.UUID,
			"label":      sshKey.Label,
			"public_key": sshKey.PublicKey,
			"comment":    sshKey.Comment,
			"created_on": sshKey.CreatedOn,
			"updated_on": sshKey.UpdatedOn,
			"links":      sshKey.Links,
		}
	}

	d.Set("ssh_keys", sshKeys)
}
