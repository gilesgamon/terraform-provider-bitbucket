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

func dataSSHKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSSHKeysRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
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
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_used": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeMap,
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

func dataSSHKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataSSHKeysRead", dumpResourceData(d, dataSSHKeys().Schema))

	url := fmt.Sprintf("2.0/workspaces/%s/ssh-keys", workspace)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from SSH keys call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate workspace %s", workspace)
	}

	if res.Body == nil {
		return diag.Errorf("error reading SSH keys with params (%s): ", dumpResourceData(d, dataSSHKeys().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	keysBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] SSH keys response: %v", keysBody)

	var keysResponse SSHKeysResponse
	decodeerr := json.Unmarshal(keysBody, &keysResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/ssh-keys", workspace))
	flattenSSHKeys(&keysResponse, d)
	return nil
}

// SSHKeysResponse represents the response from the SSH keys API
type SSHKeysResponse struct {
	Values []SSHKey `json:"values"`
	Page   int      `json:"page"`
	Size   int      `json:"size"`
	Next   string   `json:"next"`
}

// SSHKey represents an SSH key
type SSHKey struct {
	UUID      string                 `json:"uuid"`
	Key       string                 `json:"key"`
	Label     string                 `json:"label"`
	Comment   string                 `json:"comment"`
	CreatedOn string                 `json:"created_on"`
	LastUsed  string                 `json:"last_used"`
	Owner     map[string]interface{} `json:"owner"`
	Links     map[string]interface{} `json:"links"`
}

// Flattens the SSH keys information
func flattenSSHKeys(c *SSHKeysResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	keys := make([]interface{}, len(c.Values))
	for i, key := range c.Values {
		keys[i] = map[string]interface{}{
			"uuid":       key.UUID,
			"key":        key.Key,
			"label":      key.Label,
			"comment":    key.Comment,
			"created_on": key.CreatedOn,
			"last_used":  key.LastUsed,
			"owner":      key.Owner,
			"links":      key.Links,
		}
	}

	d.Set("ssh_keys", keys)
}
