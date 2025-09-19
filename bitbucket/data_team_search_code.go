package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataTeamSearchCode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTeamSearchCodeRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Team username",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"search_query": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The search query string",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number for pagination",
			},
			"pagelen": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Number of results per page",
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Result type",
						},
						"content_match_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of content matches",
						},
						"content_matches": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Content matches",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lines": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Matching lines",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"line": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Line number",
												},
												"segments": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Line segments",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"text": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Segment text",
															},
															"match": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether this segment is a match",
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
						"path_matches": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Path matches",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"text": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Path text",
									},
									"match": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether this path is a match",
									},
								},
							},
						},
						"file": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "File information",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "File path",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "File type",
									},
									"links": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "File links",
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
				},
			},
			"query_substituted": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The actual query that was executed",
			},
		},
	}
}

func dataTeamSearchCodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	username := d.Get("username").(string)
	searchQuery := d.Get("search_query").(string)
	page := d.Get("page").(int)
	pagelen := d.Get("pagelen").(int)

	// Build query parameters
	params := url.Values{}
	params.Add("search_query", searchQuery)
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("pagelen", fmt.Sprintf("%d", pagelen))

	endpoint := fmt.Sprintf("2.0/teams/%s/search/code?%s", username, params.Encode())

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from team search call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to search code for team %s", username)
	}

	if res.Body == nil {
		return diag.Errorf("error reading search results: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var searchResponse struct {
		Values           []CodeSearchResult `json:"values"`
		QuerySubstituted string             `json:"query_substituted"`
		Next             string             `json:"next"`
		Size             int                `json:"size"`
		Page             int                `json:"page"`
	}

	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return diag.FromErr(err)
	}

	var results []map[string]interface{}
	for _, result := range searchResponse.Values {
		resultMap := map[string]interface{}{
			"type":                result.Type,
			"content_match_count": result.ContentMatchCount,
		}

		// Process content matches
		var contentMatches []map[string]interface{}
		for _, match := range result.ContentMatches {
			matchMap := map[string]interface{}{}

			var lines []map[string]interface{}
			for _, line := range match.Lines {
				lineMap := map[string]interface{}{
					"line": line.Line,
				}

				var segments []map[string]interface{}
				for _, segment := range line.Segments {
					segments = append(segments, map[string]interface{}{
						"text":  segment.Text,
						"match": segment.Match,
					})
				}
				lineMap["segments"] = segments
				lines = append(lines, lineMap)
			}
			matchMap["lines"] = lines
			contentMatches = append(contentMatches, matchMap)
		}
		resultMap["content_matches"] = contentMatches

		// Process path matches
		var pathMatches []map[string]interface{}
		for _, pathMatch := range result.PathMatches {
			pathMatches = append(pathMatches, map[string]interface{}{
				"text":  pathMatch.Text,
				"match": pathMatch.Match,
			})
		}
		resultMap["path_matches"] = pathMatches

		// Process file information
		if result.File != nil {
			fileMap := map[string]interface{}{
				"path": result.File.Path,
				"type": result.File.Type,
			}

			if result.File.Links != nil {
				links := []map[string]interface{}{
					{
						"self": []map[string]interface{}{
							{
								"href": result.File.Links.Self.Href,
							},
						},
					},
				}
				fileMap["links"] = links
			}

			resultMap["file"] = []map[string]interface{}{fileMap}
		}

		results = append(results, resultMap)
	}

	d.SetId(fmt.Sprintf("team-search-%s-%s", username, searchQuery))
	d.Set("results", results)
	d.Set("query_substituted", searchResponse.QuerySubstituted)

	log.Printf("[DEBUG] Found %d search results for team %s with query: %s", len(results), username, searchQuery)

	return nil
}

