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

func dataTeamPipelineVariables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTeamPipelineVariablesRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Team username",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"variables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Variable UUID",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Variable key",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Variable value",
						},
						"secured": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the variable is secured",
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
					},
				},
			},
		},
	}
}

func dataTeamPipelineVariablesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	username := d.Get("username").(string)

	endpoint := fmt.Sprintf("2.0/teams/%s/pipelines_config/variables", username)

	res, err := client.GetAll(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from team pipeline variables call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate team %s or team pipeline variables", username)
	}

	if res.Body == nil {
		return diag.Errorf("error reading team pipeline variables: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var variablesResponse struct {
		Values []TeamPipelineVariable `json:"values"`
		Next   string                 `json:"next"`
		Size   int                    `json:"size"`
		Page   int                    `json:"page"`
	}

	if err := json.Unmarshal(body, &variablesResponse); err != nil {
		return diag.FromErr(err)
	}

	var variables []map[string]interface{}
	for _, variable := range variablesResponse.Values {
		variableMap := map[string]interface{}{
			"uuid":       variable.UUID,
			"key":        variable.Key,
			"value":      variable.Value,
			"secured":    variable.Secured,
			"created_on": variable.CreatedOn,
			"updated_on": variable.UpdatedOn,
		}
		variables = append(variables, variableMap)
	}

	d.SetId(fmt.Sprintf("team-pipeline-variables-%s", username))
	d.Set("variables", variables)

	log.Printf("[DEBUG] Found %d pipeline variables for team %s", len(variables), username)

	return nil
}

// TeamPipelineVariable represents a team pipeline variable
type TeamPipelineVariable struct {
	UUID      string `json:"uuid"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Secured   bool   `json:"secured"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}
