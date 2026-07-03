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

func dataPullRequestTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPullRequestTasksRead,
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
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}

func dataPullRequestTasksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	pullRequestID := d.Get("pull_request_id").(string)

	endpoint := fmt.Sprintf("2.0/repositories/%s/%s/pullrequests/%s/tasks", workspace, repoSlug, pullRequestID)

	res, err := client.GetAll(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from pull request tasks call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s or pull request %s", workspace, repoSlug, pullRequestID)
	}

	if res.Body == nil {
		return diag.Errorf("error reading pull request tasks: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var tasksResponse struct {
		Values []PullRequestTask `json:"values"`
		Next   string            `json:"next"`
		Size   int               `json:"size"`
		Page   int               `json:"page"`
	}

	if err := json.Unmarshal(body, &tasksResponse); err != nil {
		return diag.FromErr(err)
	}

	var tasks []map[string]interface{}
	for _, task := range tasksResponse.Values {
		taskMap := map[string]interface{}{
			"id":         task.ID,
			"content":    task.Content,
			"state":      task.State,
			"created_on": task.CreatedOn,
			"updated_on": task.UpdatedOn,
		}

		if task.Creator != nil {
			taskMap["creator"] = []map[string]interface{}{
				{
					"display_name": task.Creator.DisplayName,
					"uuid":         task.Creator.UUID,
					"username":     task.Creator.Username,
				},
			}
		}

		tasks = append(tasks, taskMap)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", workspace, repoSlug, pullRequestID))
	d.Set("tasks", tasks)

	log.Printf("[DEBUG] Found %d tasks for pull request %s in repository %s/%s", len(tasks), pullRequestID, workspace, repoSlug)

	return nil
}

// PullRequestTask represents a pull request task
type PullRequestTask struct {
	ID        int      `json:"id"`
	Content   string   `json:"content"`
	State     string   `json:"state"`
	CreatedOn string   `json:"created_on"`
	UpdatedOn string   `json:"updated_on"`
	Creator   *Account `json:"creator,omitempty"`
}
