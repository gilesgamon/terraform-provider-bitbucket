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

func dataIssueTransitions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIssueTransitionsRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issue_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"transitions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"to": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fields": {
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
				},
			},
		},
	}
}

func dataIssueTransitionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	issueID := d.Get("issue_id").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataIssueTransitionsRead", dumpResourceData(d, dataIssueTransitions().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/issues/%s/transitions", workspace, repoSlug, issueID)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from issue transitions call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate issue %s in repository %s/%s", issueID, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading issue transitions with params (%s): ", dumpResourceData(d, dataIssueTransitions().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	transitionsBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] issue transitions response: %v", transitionsBody)

	var transitionsResponse IssueTransitionsResponse
	decodeerr := json.Unmarshal(transitionsBody, &transitionsResponse)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/issues/%s/transitions", workspace, repoSlug, issueID))
	flattenIssueTransitions(&transitionsResponse, d)
	return nil
}

// IssueTransitionsResponse represents the response from the issue transitions API
type IssueTransitionsResponse struct {
	Values []IssueTransition `json:"values"`
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Next   string            `json:"next"`
}

// IssueTransition represents an issue transition
type IssueTransition struct {
	UUID   string                 `json:"uuid"`
	Name   string                 `json:"name"`
	To     map[string]interface{} `json:"to"`
	Fields map[string]interface{} `json:"fields"`
	Links  map[string]interface{} `json:"links"`
}

// Flattens the issue transitions information
func flattenIssueTransitions(c *IssueTransitionsResponse, d *schema.ResourceData) {
	if c == nil {
		return
	}

	transitions := make([]interface{}, len(c.Values))
	for i, transition := range c.Values {
		transitions[i] = map[string]interface{}{
			"uuid":   transition.UUID,
			"name":   transition.Name,
			"to":     transition.To,
			"fields": transition.Fields,
			"links":  transition.Links,
		}
	}

	d.Set("transitions", transitions)
}
