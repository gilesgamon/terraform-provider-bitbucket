package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProjectDeployKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataProjectDeployKeysRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deploy_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"added_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_used": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataProjectDeployKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	projectKey := d.Get("project_key").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataProjectDeployKeysRead", dumpResourceData(d, dataProjectDeployKeys().Schema))

	endpoint := fmt.Sprintf("2.0/workspaces/%s/projects/%s/deploy-keys", workspace, projectKey)

	client := m.(Clients).httpClient

	var deployKeys []interface{}
	for endpoint != "" {
		res, err := client.Get(endpoint)
		if err != nil {
			return diag.FromErr(err)
		}

		if res.StatusCode == http.StatusNotFound {
			return diag.Errorf("unable to locate project %s in workspace %s", projectKey, workspace)
		}

		if err := handleClientError(res, err); err != nil {
			return diag.FromErr(err)
		}

		body, readerr := io.ReadAll(res.Body)
		res.Body.Close()
		if readerr != nil {
			return diag.FromErr(readerr)
		}

		var page struct {
			Values []ProjectDeployKey `json:"values"`
			Next   string             `json:"next"`
		}
		if decodeerr := json.Unmarshal(body, &page); decodeerr != nil {
			return diag.FromErr(decodeerr)
		}

		for _, deployKey := range page.Values {
			deployKeys = append(deployKeys, map[string]interface{}{
				"key_id":    fmt.Sprintf("%d", deployKey.ID),
				"key":       deployKey.Key,
				"label":     deployKey.Label,
				"comment":   deployKey.Comment,
				"added_on":  deployKey.AddedOn,
				"last_used": deployKey.LastUsed,
			})
		}

		endpoint = nextPageEndpoint(page.Next)
	}

	d.SetId(fmt.Sprintf("%s/%s/deploy-keys", workspace, projectKey))
	d.Set("deploy_keys", deployKeys)
	return nil
}

// nextPageEndpoint converts an absolute Bitbucket `next` URL into an endpoint
// relative to the API base, which is what the HTTP client expects.
func nextPageEndpoint(raw string) string {
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	rel := strings.TrimLeft(parsed.Path, "/")
	if parsed.RawQuery != "" {
		rel += "?" + parsed.RawQuery
	}
	return rel
}
