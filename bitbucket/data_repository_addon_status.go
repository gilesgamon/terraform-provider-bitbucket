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

func dataRepositoryAddonStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryAddonStatusRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addon_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"details": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataRepositoryAddonStatusRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	addonKey := d.Get("addon_key").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryAddonStatusRead", dumpResourceData(d, dataRepositoryAddonStatus().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/addons/%s/status", workspace, repoSlug, addonKey)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository addon status call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate addon %s for repository %s/%s", addonKey, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository addon status with params (%s): ", dumpResourceData(d, dataRepositoryAddonStatus().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	statusBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository addon status response: %v", statusBody)

	var addonStatus RepositoryAddonStatus
	decodeerr := json.Unmarshal(statusBody, &addonStatus)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/addons/%s/status", workspace, repoSlug, addonKey))
	flattenRepositoryAddonStatus(&addonStatus, d)
	return nil
}

// RepositoryAddonStatus represents the response from the repository addon status API
type RepositoryAddonStatus struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Flattens the repository addon status information
func flattenRepositoryAddonStatus(c *RepositoryAddonStatus, d *schema.ResourceData) {
	if c == nil {
		return
	}

	d.Set("status", c.Status)
	d.Set("message", c.Message)
	d.Set("details", c.Details)
}
