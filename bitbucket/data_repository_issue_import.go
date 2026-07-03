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

func dataRepositoryIssueImport() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryIssueImportRead,
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
			"import_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Import status",
			},
			"import_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Import URL",
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

func dataRepositoryIssueImportRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	endpoint := fmt.Sprintf("2.0/repositories/%s/%s/issues/import", workspace, repoSlug)

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from issue import call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s or issue import", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue import: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var importResponse IssueImport
	if err := json.Unmarshal(body, &importResponse); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", workspace, repoSlug))
	d.Set("import_status", importResponse.Status)
	d.Set("import_url", importResponse.ImportURL)
	d.Set("created_on", importResponse.CreatedOn)
	d.Set("updated_on", importResponse.UpdatedOn)

	log.Printf("[DEBUG] Retrieved issue import status: %s for repository %s/%s", importResponse.Status, workspace, repoSlug)

	return nil
}

// IssueImport represents an issue import
type IssueImport struct {
	Status    string `json:"status"`
	ImportURL string `json:"import_url"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}
