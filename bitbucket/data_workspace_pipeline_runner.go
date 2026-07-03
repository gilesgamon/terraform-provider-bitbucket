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

func dataWorkspacePipelineRunner() *schema.Resource {
	runnerSchema := pipelineRunnerSchema()
	runnerSchema["workspace"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	runnerSchema["runner_uuid"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}

	return &schema.Resource{
		ReadContext: dataWorkspacePipelineRunnerRead,
		Schema:      runnerSchema,
	}
}

func dataWorkspacePipelineRunnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	runnerUUID := d.Get("runner_uuid").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataWorkspacePipelineRunnerRead", dumpResourceData(d, dataWorkspacePipelineRunner().Schema))

	url := fmt.Sprintf("2.0/workspaces/%s/pipelines-config/runners/%s", workspace, runnerUUID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from workspace pipeline runner call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate runner %s in workspace %s", runnerUUID, workspace)
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	runnerBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)

	var runner PipelineRunner
	decodeerr := json.Unmarshal(runnerBody, &runner)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/pipelines-config/runners/%s", workspace, runnerUUID))
	flattened := flattenPipelineRunner(&runner)
	for k, v := range flattened {
		d.Set(k, v)
	}
	return nil
}
