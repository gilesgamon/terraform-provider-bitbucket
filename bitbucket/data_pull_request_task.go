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

func dataPullRequestTask() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPullRequestTaskRead,
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
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Task ID",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task content",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task state",
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
			"creator": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Task creator",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataPullRequestTaskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pullRequestID := d.Get("pull_request_id").(string)
	taskID := d.Get("task_id").(string)

	endpoint := fmt.Sprintf("2.0/repositories/%s/%s/pullrequests/%s/tasks/%s", workspace, repoSlug, pullRequestID, taskID)

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from pull request task call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s, pull request %s, or task %s", workspace, repoSlug, pullRequestID, taskID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pull request task: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var task PullRequestTask
	if err := json.Unmarshal(body, &task); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", workspace, repoSlug, pullRequestID, taskID))
	d.Set("id", task.ID)
	d.Set("content", task.Content)
	d.Set("state", task.State)
	d.Set("created_on", task.CreatedOn)
	d.Set("updated_on", task.UpdatedOn)

	if task.Creator != nil {
		creator := []map[string]interface{}{
			{
				"display_name": task.Creator.DisplayName,
				"uuid":         task.Creator.UUID,
				"username":     task.Creator.Username,
			},
		}
		d.Set("creator", creator)
	}

	log.Printf("[DEBUG] Retrieved pull request task: %d for pull request %s in repository %s/%s", task.ID, pullRequestID, workspace, repoSlug)

	return nil
}
