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

func dataRepositoryIssueExport() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryIssueExportRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"export_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"export": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
						"download_url": {
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

func dataRepositoryIssueExportRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	exportID := d.Get("export_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryIssueExportRead", dumpResourceData(d, dataRepositoryIssueExport().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues/export/%s", workspace, repoSlug, exportID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository issue export call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate export %s for repository %s/%s", exportID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository issue export with params (%s): ", dumpResourceData(d, dataRepositoryIssueExport().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	exportBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository issue export response: %v", exportBody)

	var issueExport RepositoryIssueExport
	decodeerr := json.Unmarshal(exportBody, &issueExport)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues/export/%s", workspace, repoSlug, exportID))
	flattenRepositoryIssueExport(&issueExport, d)
	return nil
}

// RepositoryIssueExport represents the response from the repository issue export API
type RepositoryIssueExport struct {
	UUID         string                 `json:"uuid"`
	Status       string                 `json:"status"`
	CreatedOn    string                 `json:"created_on"`
	UpdatedOn    string                 `json:"updated_on"`
	DownloadURL  string                 `json:"download_url"`
	Links        map[string]interface{} `json:"links"`
}

// Flattens the repository issue export information
func flattenRepositoryIssueExport(c *RepositoryIssueExport, d *schema.ResourceData) {
	if c == nil {
		return
	}

	export := map[string]interface{}{
		"uuid":          c.UUID,
		"status":        c.Status,
		"created_on":    c.CreatedOn,
		"updated_on":    c.UpdatedOn,
		"download_url":  c.DownloadURL,
		"links":         c.Links,
	}

	d.Set("export", []interface{}{export})
}
