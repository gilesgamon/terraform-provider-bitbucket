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

func dataRepositorySettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositorySettingsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"fork_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"language": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_issues": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_wiki": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"website": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"mainbranch": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"links": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataRepositorySettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataRepositorySettingsRead", dumpResourceData(d, dataRepositorySettings().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s", workspace, repoSlug)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from repository settings call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository settings with params (%s): ", dumpResourceData(d, dataRepositorySettings().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	settingsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] repository settings response: %v", settingsBody)

	var repository RepositorySettings
	decodeerr := json.Unmarshal(settingsBody, &repository)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s", workspace, repoSlug))
	flattenRepositorySettings(&repository, d)
	return nil
}

// RepositorySettings represents repository settings from the API
type RepositorySettings struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	IsPrivate   bool                   `json:"is_private"`
	ForkPolicy  string                 `json:"fork_policy"`
	Language    string                 `json:"language"`
	HasIssues   bool                   `json:"has_issues"`
	HasWiki     bool                   `json:"has_wiki"`
	Size        int                    `json:"size"`
	UpdatedOn   string                 `json:"updated_on"`
	CreatedOn   string                 `json:"created_on"`
	Scm         string                 `json:"scm"`
	Website     string                 `json:"website"`
	Project     map[string]interface{} `json:"project"`
	Mainbranch  map[string]interface{} `json:"mainbranch"`
	Links       map[string]interface{} `json:"links"`
}

// Flattens the repository settings information
func flattenRepositorySettings(r *RepositorySettings, d *schema.ResourceData) {
	if r == nil {
		return
	}

	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("is_private", r.IsPrivate)
	d.Set("fork_policy", r.ForkPolicy)
	d.Set("language", r.Language)
	d.Set("has_issues", r.HasIssues)
	d.Set("has_wiki", r.HasWiki)
	d.Set("size", r.Size)
	d.Set("updated_on", r.UpdatedOn)
	d.Set("created_on", r.CreatedOn)
	d.Set("scm", r.Scm)
	d.Set("website", r.Website)
	d.Set("project", r.Project)
	d.Set("mainbranch", r.Mainbranch)
	d.Set("links", r.Links)
}
