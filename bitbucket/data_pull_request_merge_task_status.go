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

func dataPullRequestMergeTaskStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPullRequestMergeTaskStatusRead,
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
			"pull_request_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Pull request ID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"task_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Task ID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Merge task status",
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

func dataPullRequestMergeTaskStatusRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pullRequestID := d.Get("pull_request_id").(string)
	taskID := d.Get("task_id").(string)

	endpoint := fmt.Sprintf("2.0/repositories/%s/%s/pullrequests/%s/merge/task-status/%s", workspace, repoSlug, pullRequestID, taskID)

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from pull request merge task status call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s, pull request %s, or merge task %s", workspace, repoSlug, pullRequestID, taskID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pull request merge task status: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var taskStatus MergeTaskStatus
	if err := json.Unmarshal(body, &taskStatus); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", workspace, repoSlug, pullRequestID, taskID))
	d.Set("status", taskStatus.Status)
	d.Set("created_on", taskStatus.CreatedOn)
	d.Set("updated_on", taskStatus.UpdatedOn)

	log.Printf("[DEBUG] Retrieved merge task status: %s for pull request %s in repository %s/%s", taskStatus.Status, pullRequestID, workspace, repoSlug)

	return nil
}

// MergeTaskStatus represents a merge task status
type MergeTaskStatus struct {
	Status    string `json:"status"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}
