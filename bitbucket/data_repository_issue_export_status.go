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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataRepositoryIssueExportStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryIssueExportStatusRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Workspace slug or UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"repo_slug": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Repository slug or UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"export_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Export status",
			},
			"export_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Export URL",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp",
			},
			"updated_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp",
			},
		},
	}
}

func dataRepositoryIssueExportStatusRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	endpoint := fmt.Sprintf("2.0/repositories/%s/%s/issues/export", workspace, repoSlug)

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from issue export call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s or issue export", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue export: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var exportResponse IssueExport
	if err := json.Unmarshal(body, &exportResponse); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", workspace, repoSlug))
	d.Set("export_status", exportResponse.Status)
	d.Set("export_url", exportResponse.ExportURL)
	d.Set("created_on", exportResponse.CreatedOn)
	d.Set("updated_on", exportResponse.UpdatedOn)

	log.Printf("[DEBUG] Retrieved issue export status: %s for repository %s/%s", exportResponse.Status, workspace, repoSlug)

	return nil
}

// IssueExport represents an issue export
type IssueExport struct {
	Status    string `json:"status"`
	ExportURL string `json:"export_url"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}

