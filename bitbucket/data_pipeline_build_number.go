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

func dataPipelineBuildNumber() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPipelineBuildNumberRead,
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
			"next": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Next build number",
			},
		},
	}
}

func dataPipelineBuildNumberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	endpoint := fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/build_number", workspace, repoSlug)

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from pipeline build number call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s or pipeline build number", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pipeline build number: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var buildNumberResponse struct {
		Next int `json:"next"`
	}

	if err := json.Unmarshal(body, &buildNumberResponse); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", workspace, repoSlug))
	d.Set("next", buildNumberResponse.Next)

	log.Printf("[DEBUG] Retrieved next build number: %d for repository %s/%s", buildNumberResponse.Next, workspace, repoSlug)

	return nil
}

