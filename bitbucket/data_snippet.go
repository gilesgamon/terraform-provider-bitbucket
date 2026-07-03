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

func dataSnippet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSnippetRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Workspace slug or UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"encoded_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Snippet encoded ID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
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
	}
}

func dataSnippetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	encodedID := d.Get("encoded_id").(string)

	endpoint := fmt.Sprintf("2.0/snippets/%s/%s", workspace, encodedID)

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from snippet call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate snippet %s in workspace %s", encodedID, workspace)
	}

	if res.Body == nil {
		return diag.Errorf("error reading snippet: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var snippet Snippet
	if err := json.Unmarshal(body, &snippet); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", workspace, encodedID))
	d.Set("id", snippet.ID)
	d.Set("title", snippet.Title)
	d.Set("scm", snippet.Scm)
	d.Set("created_on", snippet.CreatedOn)
	d.Set("updated_on", snippet.UpdatedOn)
	d.Set("is_private", snippet.IsPrivate)

	if snippet.Owner != nil {
		owner := []map[string]interface{}{
			{
				"display_name": snippet.Owner.DisplayName,
				"uuid":         snippet.Owner.UUID,
				"username":     snippet.Owner.Username,
			},
		}
		d.Set("owner", owner)
	}

	if snippet.Creator != nil {
		creator := []map[string]interface{}{
			{
				"display_name": snippet.Creator.DisplayName,
				"uuid":         snippet.Creator.UUID,
				"username":     snippet.Creator.Username,
			},
		}
		d.Set("creator", creator)
	}

	if snippet.Links != nil {
		links := []map[string]interface{}{
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
		d.Set("links", links)
	}

	log.Printf("[DEBUG] Retrieved snippet: %s", snippet.Title)

	return nil
}
