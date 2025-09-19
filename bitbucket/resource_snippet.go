package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSnippet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnippetCreate,
		ReadContext:   resourceSnippetRead,
		UpdateContext: resourceSnippetUpdate,
		DeleteContext: resourceSnippetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Workspace slug or UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"title": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Snippet title",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"scm": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "git",
				Description:  "The DVCS used to store the snippet",
				ValidateFunc: validation.StringInSlice([]string{"git"}, false),
			},
			"is_private": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the snippet is private",
			},
			"files": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "Snippet files (filename -> content)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Snippet ID",
			},
			"encoded_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Snippet encoded ID",
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

func resourceSnippetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	title := d.Get("title").(string)
	scm := d.Get("scm").(string)
	isPrivate := d.Get("is_private").(bool)
	files := d.Get("files").(map[string]interface{})

	// Convert files to the expected format
	snippetFiles := make(map[string]SnippetFile)
	for filename, content := range files {
		snippetFiles[filename] = SnippetFile{
			Content: content.(string),
		}
	}

	snippetRequest := SnippetRequest{
		Title:     title,
		Scm:       scm,
		IsPrivate: isPrivate,
		Files:     snippetFiles,
	}

	jsonPayload, err := json.Marshal(snippetRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	endpoint := fmt.Sprintf("2.0/snippets/%s", workspace)
	res, err := client.Post(endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from snippet creation")
	}

	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		return diag.Errorf("failed to create snippet: %s", string(body))
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

	// Extract encoded_id from the response
	encodedID := extractEncodedIDFromSnippet(snippet)
	if encodedID == "" {
		return diag.Errorf("failed to extract encoded_id from snippet response")
	}

	d.SetId(fmt.Sprintf("%s/%s", workspace, encodedID))

	log.Printf("[DEBUG] Created snippet: %s with ID: %s", title, d.Id())

	return resourceSnippetRead(ctx, d, m)
}

func resourceSnippetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace, encodedID, err := snippetId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	endpoint := fmt.Sprintf("2.0/snippets/%s/%s", workspace, encodedID)
	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from snippet call")
	}

	if res.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] Snippet (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
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

	d.Set("workspace", workspace)
	d.Set("encoded_id", encodedID)
	d.Set("id", fmt.Sprintf("%d", snippet.ID))
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

func resourceSnippetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace, encodedID, err := snippetId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	title := d.Get("title").(string)
	scm := d.Get("scm").(string)
	isPrivate := d.Get("is_private").(bool)
	files := d.Get("files").(map[string]interface{})

	// Convert files to the expected format
	snippetFiles := make(map[string]SnippetFile)
	for filename, content := range files {
		snippetFiles[filename] = SnippetFile{
			Content: content.(string),
		}
	}

	snippetRequest := SnippetRequest{
		Title:     title,
		Scm:       scm,
		IsPrivate: isPrivate,
		Files:     snippetFiles,
	}

	jsonPayload, err := json.Marshal(snippetRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	endpoint := fmt.Sprintf("2.0/snippets/%s/%s", workspace, encodedID)
	res, err := client.Put(endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from snippet update")
	}

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return diag.Errorf("failed to update snippet: %s", string(body))
	}

	log.Printf("[DEBUG] Updated snippet: %s", title)

	return resourceSnippetRead(ctx, d, m)
}

func resourceSnippetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace, encodedID, err := snippetId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	endpoint := fmt.Sprintf("2.0/snippets/%s/%s", workspace, encodedID)
	res, err := client.Delete(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from snippet deletion")
	}

	if res.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(res.Body)
		return diag.Errorf("failed to delete snippet: %s", string(body))
	}

	log.Printf("[DEBUG] Deleted snippet: %s", d.Id())

	return nil
}

// Helper functions
func snippetId(id string) (workspace, encodedID string, err error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected format of ID (%q), expected WORKSPACE/ENCODED-ID", id)
	}
	return parts[0], parts[1], nil
}

func extractEncodedIDFromSnippet(snippet Snippet) string {
	// Extract encoded_id from the snippet's self link
	if snippet.Links != nil && snippet.Links.Self.Href != "" {
		// Parse the href to extract the encoded_id
		// Example: https://api.bitbucket.org/2.0/snippets/workspace/encoded_id
		href := snippet.Links.Self.Href
		parts := strings.Split(href, "/")
		if len(parts) > 0 {
			// Get the last part which should be the encoded_id
			lastPart := parts[len(parts)-1]
			return lastPart
		}
	}
	return ""
}

// SnippetRequest represents the request payload for creating/updating snippets
type SnippetRequest struct {
	Title     string                 `json:"title"`
	Scm       string                 `json:"scm"`
	IsPrivate bool                   `json:"is_private"`
	Files     map[string]SnippetFile `json:"files"`
}

// SnippetFile represents a file in a snippet
type SnippetFile struct {
	Content string `json:"content"`
}
