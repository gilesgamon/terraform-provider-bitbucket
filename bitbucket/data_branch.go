package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/DrFaust92/bitbucket-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataBranch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataBranchRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"branch_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Branch name (e.g., 'main', 'develop', 'feature/new-feature')",
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_author": {
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
					},
				},
			},
			"target_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataBranchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace := d.Get("workspace").(string)
	repoSlug := d.Get("repo_slug").(string)
	branchName := d.Get("branch_name").(string)

	log.Printf("[DEBUG]: params for %s: %v", "dataBranchRead", dumpResourceData(d, dataBranch().Schema))

	url := fmt.Sprintf("2.0/repositories/%s/%s/refs/branches/%s",
		workspace,
		repoSlug,
		branchName,
	)

	client := m.(Clients).httpClient
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	if res == nil {
		return diag.Errorf("no response returned from branch call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate branch %s in repository %s/%s", branchName, workspace, repoSlug)
	}

	if res.Body == nil {
		return diag.Errorf("error reading branch information with params (%s): ", dumpResourceData(d, dataBranch().Schema))
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	branchBody, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}
	log.Printf("[DEBUG] http response: %v", res)
	log.Printf("[DEBUG] branch response: %v", branchBody)

	var branch Branch
	decodeerr := json.Unmarshal(branchBody, &branch)
	if decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", workspace, repoSlug, branch.Name))
	flattenBranch(&branch, d)
	return nil
}

// Branch represents a Bitbucket branch
type Branch struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Target     BranchTarget           `json:"target"`
	Repository bitbucket.Repository   `json:"repository"`
	Links      map[string]interface{} `json:"links"`
}

// BranchTarget represents the target of a branch (usually a commit)
type BranchTarget struct {
	Hash  string                 `json:"hash"`
	Type  string                 `json:"type"`
	Links map[string]interface{} `json:"links"`
}

// Flattens the branch information
func flattenBranch(b *Branch, d *schema.ResourceData) {
	if b == nil {
		return
	}

	d.Set("name", b.Name)
	d.Set("target_hash", b.Target.Hash)
	d.Set("target_date", b.Target.Type)
	d.Set("target_message", b.Type)
}
