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

func dataSnippets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSnippetsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Workspace slug or UUID to filter snippets",
			},
			"snippets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Snippet ID",
						},
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Snippet title",
						},
						"scm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The DVCS used to store the snippet",
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
						"owner": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Snippet owner",
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
						"creator": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Snippet creator",
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
						"is_private": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the snippet is private",
						},
						"links": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Snippet links",
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
									"html": {
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

func dataSnippetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)

	var endpoint string
	if workspace != "" {
		endpoint = fmt.Sprintf("2.0/snippets/%s", workspace)
	} else {
		endpoint = "2.0/snippets"
	}

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from snippets call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate snippets for workspace %s", workspace)
	}

	if res.Body == nil {
		return diag.Errorf("error reading snippets: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var snippetsResponse struct {
		Values []Snippet `json:"values"`
		Next   string    `json:"next"`
		Size   int       `json:"size"`
		Page   int       `json:"page"`
	}

	if err := json.Unmarshal(body, &snippetsResponse); err != nil {
		return diag.FromErr(err)
	}

	var snippets []map[string]interface{}
	for _, snippet := range snippetsResponse.Values {
		snippetMap := map[string]interface{}{
			"id":         snippet.ID,
			"title":      snippet.Title,
			"scm":        snippet.Scm,
			"created_on": snippet.CreatedOn,
			"updated_on": snippet.UpdatedOn,
			"is_private": snippet.IsPrivate,
		}

		if snippet.Owner != nil {
			snippetMap["owner"] = []map[string]interface{}{
				{
					"display_name": snippet.Owner.DisplayName,
					"uuid":         snippet.Owner.UUID,
					"username":     snippet.Owner.Username,
				},
			}
		}

		if snippet.Creator != nil {
			snippetMap["creator"] = []map[string]interface{}{
				{
					"display_name": snippet.Creator.DisplayName,
					"uuid":         snippet.Creator.UUID,
					"username":     snippet.Creator.Username,
				},
			}
		}

		if snippet.Links != nil {
			snippetMap["links"] = []map[string]interface{}{
				{
					"self": []map[string]interface{}{
						{
							"href": snippet.Links.Self.Href,
						},
					},
					"html": []map[string]interface{}{
						{
							"href": snippet.Links.Html.Href,
						},
					},
				},
			}
		}

		snippets = append(snippets, snippetMap)
	}

	d.SetId(fmt.Sprintf("snippets-%s", workspace))
	d.Set("snippets", snippets)

	log.Printf("[DEBUG] Found %d snippets", len(snippets))

	return nil
}

// Snippet represents a Bitbucket snippet
type Snippet struct {
	ID        int           `json:"id"`
	Title     string        `json:"title"`
	Scm       string        `json:"scm"`
	CreatedOn string        `json:"created_on"`
	UpdatedOn string        `json:"updated_on"`
	Owner     *Account      `json:"owner,omitempty"`
	Creator   *Account      `json:"creator,omitempty"`
	IsPrivate bool          `json:"is_private"`
	Links     *SnippetLinks `json:"links,omitempty"`
}

// Account represents a Bitbucket account/user
type Account struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	UUID        string `json:"uuid"`
}

// Link represents a URL link
type Link struct {
	Href string `json:"href"`
}

// SnippetLinks represents snippet links
type SnippetLinks struct {
	Self Link `json:"self"`
	Html Link `json:"html"`
}
