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

func dataRepositoryPipelineRunners() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryPipelineRunnersRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"runners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: pipelineRunnerSchema(),
				},
			},
		},
	}
}

func dataRepositoryPipelineRunnersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryPipelineRunnersRead", dumpResourceData(d, dataRepositoryPipelineRunners().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/pipelines-config/runners", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository pipeline runners call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	runnersBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)

	var runnersResponse PipelineRunnersResponse
	decodeerr := json.Unmarshal(runnersBody, &runnersResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/pipelines-config/runners", workspace, repoSlug))
	d.Set("runners", flattenPipelineRunners(runnersResponse.Values))
	return nil
}
