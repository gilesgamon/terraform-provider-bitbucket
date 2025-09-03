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

func dataPipelineKnownHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineKnownHostsRead,
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
						"key_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fingerprint": {
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

func dataPipelineKnownHostsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataPipelineKnownHostsRead", dumpResourceData(d, dataPipelineKnownHosts().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines/ssh/known-hosts", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from pipeline known hosts call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline known hosts with params (%s): ", dumpResourceData(d, dataPipelineKnownHosts().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	knownHostsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] pipeline known hosts response: %v", knownHostsBody)

	var knownHostsResponse PipelineKnownHostsResponse
	decodeerr := json.Unmarshal(knownHostsBody, &knownHostsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines/ssh/known-hosts", workspace, repoSlug))
	flattenPipelineKnownHosts(&knownHostsResponse, d)
	return nil
}

// PipelineKnownHostsResponse represents the response from the pipeline known hosts API
type PipelineKnownHostsResponse struct {
	Values []PipelineKnownHost `json:"values"`
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	Next   string              `json:"next"`
}

// PipelineKnownHost represents a pipeline known host
type PipelineKnownHost struct {
	UUID       string                 `json:"uuid"`
	Hostname   string                 `json:"hostname"`
	PublicKey  string                 `json:"public_key"`
	KeyType    string                 `json:"key_type"`
	Fingerprint string                `json:"fingerprint"`
	CreatedOn  string                 `json:"created_on"`
	UpdatedOn  string                 `json:"updated_on"`
	Links      map[string]interface{} `json:"links"`
}

// Flattens the pipeline known hosts information
func flattenPipelineKnownHosts(c *PipelineKnownHostsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	knownHosts := make([]interface{}, len(c.Values))
	for i, knownHost := range c.Values {
		knownHosts[i] = map[string]interface{}{
			"uuid":        knownHost.UUID,
			"hostname":    knownHost.Hostname,
			"public_key":  knownHost.PublicKey,
			"key_type":    knownHost.KeyType,
			"fingerprint": knownHost.Fingerprint,
			"created_on":  knownHost.CreatedOn,
			"updated_on":  knownHost.UpdatedOn,
			"links":       knownHost.Links,
		}
	}

	d.Set("known_hosts", knownHosts)
}
