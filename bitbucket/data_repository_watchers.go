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

func dataRepositoryWatchers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryWatchersRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"watchers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nickname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_staff": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"account_status": {
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

func dataRepositoryWatchersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositoryWatchersRead", dumpResourceData(d, dataRepositoryWatchers().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/watchers", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository watchers call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository watchers with params (%s): ", dumpResourceData(d, dataRepositoryWatchers().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	watchersBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository watchers response: %v", watchersBody)

	var watchersResponse RepositoryWatchersResponse
	decodeerr := json.Unmarshal(watchersBody, &watchersResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/watchers", workspace, repoSlug))
	flattenRepositoryWatchers(&watchersResponse, d)
	return nil
}

// RepositoryWatchersResponse represents the response from the repository watchers API
type RepositoryWatchersResponse struct {
	Values []RepositoryWatcher `json:"values"`
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	Next   string              `json:"next"`
}

// RepositoryWatcher represents a repository watcher
type RepositoryWatcher struct {
	Username      string                 `json:"username"`
	DisplayName   string                 `json:"display_name"`
	UUID          string                 `json:"uuid"`
	Type          string                 `json:"type"`
	Nickname      string                 `json:"nickname"`
	AccountID     string                 `json:"account_id"`
	CreatedOn     string                 `json:"created_on"`
	IsStaff       bool                   `json:"is_staff"`
	AccountStatus string                 `json:"account_status"`
	Links         map[string]interface{} `json:"links"`
}

// Flattens the repository watchers information
func flattenRepositoryWatchers(c *RepositoryWatchersResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	watchers := make([]interface{}, len(c.Values))
	for i, watcher := range c.Values {
		watchers[i] = map[string]interface{}{
			"username":       watcher.Username,
			"display_name":   watcher.DisplayName,
			"uuid":           watcher.UUID,
			"type":           watcher.Type,
			"nickname":       watcher.Nickname,
			"account_id":     watcher.AccountID,
			"created_on":     watcher.CreatedOn,
			"is_staff":       watcher.IsStaff,
			"account_status": watcher.AccountStatus,
			"links":          watcher.Links,
		}
	}

	d.Set("watchers", watchers)
}
