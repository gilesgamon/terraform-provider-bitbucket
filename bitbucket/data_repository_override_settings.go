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

func dataRepositoryOverrideSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRepositoryOverrideSettingsRead,
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
			"settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Settings type",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Settings name",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Settings value",
						},
						"links": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Settings links",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"self": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"href": {
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
				},
			},
		},
	}
}

func dataRepositoryOverrideSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)

	endpoint := fmt.Sprintf("2.0/repositories/%s/%s/override-settings", workspace, repoSlug)

	res, err := client.GetAll(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from repository override settings call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate repository %s/%s or override settings", workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading repository override settings: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var settingsResponse struct {
		Values []OverrideSetting `json:"values"`
		Next   string            `json:"next"`
		Size   int               `json:"size"`
		Page   int               `json:"page"`
	}

	if err := json.Unmarshal(body, &settingsResponse); err != nil {
		return diag.FromErr(err)
	}

	var settings []map[string]interface{}
	for _, setting := range settingsResponse.Values {
		settingMap := map[string]interface{}{
			"type":  setting.Type,
			"name":  setting.Name,
			"value": setting.Value,
		}

		if setting.Links != nil {
			settingMap["links"] = []map[string]interface{}{
				{
					"self": []map[string]interface{}{
						{
							"href": setting.Links.Self.Href,
						},
					},
				},
			}
		}

		settings = append(settings, settingMap)
	}

	d.SetId(fmt.Sprintf("%s/%s", workspace, repoSlug))
	d.Set("settings", settings)

	log.Printf("[DEBUG] Found %d override settings for repository %s/%s", len(settings), workspace, repoSlug)

	return nil
}

// OverrideSetting represents a repository override setting
type OverrideSetting struct {
	Type  string                `json:"type"`
	Name  string                `json:"name"`
	Value string                `json:"value"`
	Links *OverrideSettingLinks `json:"links,omitempty"`
}

// OverrideSettingLinks represents override setting links
type OverrideSettingLinks struct {
	Self Link `json:"self"`
}
