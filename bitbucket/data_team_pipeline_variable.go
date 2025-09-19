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

func dataTeamPipelineVariable() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTeamPipelineVariableRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Team username",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"variable_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Variable UUID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
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
	}
}

func dataTeamPipelineVariableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	username := d.Get("username").(string)
	variableUUID := d.Get("variable_uuid").(string)

	endpoint := fmt.Sprintf("2.0/teams/%s/pipelines_config/variables/%s", username, variableUUID)

	res, err := client.Get(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("no response returned from team pipeline variable call. Make sure your credentials are accurate.")
	}

	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("unable to locate team pipeline variable %s for team %s", variableUUID, username)
	}

	if res.Body == nil {
		return diag.Errorf("error reading team pipeline variable: empty response body")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var variable TeamPipelineVariable
	if err := json.Unmarshal(body, &variable); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", username, variableUUID))
	d.Set("uuid", variable.UUID)
	d.Set("key", variable.Key)
	d.Set("value", variable.Value)
	d.Set("secured", variable.Secured)
	d.Set("created_on", variable.CreatedOn)
	d.Set("updated_on", variable.UpdatedOn)

	log.Printf("[DEBUG] Retrieved team pipeline variable: %s for team %s", variable.Key, username)

	return nil
}
