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
)

// ProjectDeployKey represents a deploy (access) key configured at the project
// level. Project deploy keys are inherited by all repositories in the project.
type ProjectDeployKey struct {
	ID       int    `json:"id,omitempty"`
	Key      string `json:"key,omitempty"`
	Label    string `json:"label,omitempty"`
	Comment  string `json:"comment,omitempty"`
	AddedOn  string `json:"added_on,omitempty"`
	LastUsed string `json:"last_used,omitempty"`
}

func resourceProjectDeployKey() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceProjectDeployKeyCreate,
		ReadWithoutTimeout:   resourceProjectDeployKeyRead,
		DeleteWithoutTimeout: resourceProjectDeployKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The workspace ID (slug) or the workspace UUID surrounded by curly-braces.",
			},
			"project_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The project key (for example `PROJ`).",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The public SSH key value.",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The user-defined label for the deploy key.",
			},
			"key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"added_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_used": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func expandProjectDeployKey(d *schema.ResourceData) *ProjectDeployKey {
	key := &ProjectDeployKey{
		Key: d.Get("key").(string),
	}

	if v, ok := d.GetOk("label"); ok {
		key.Label = v.(string)
	}

	return key
}

func resourceProjectDeployKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	projectKey := d.Get("project_key").(string)

	payload, err := json.Marshal(expandProjectDeployKey(d))
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("2.0/workspaces/%s/projects/%s/deploy-keys", workspace, projectKey)
	res, err := client.Post(url, bytes.NewBuffer(payload))
	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	body, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}

	log.Printf("[DEBUG] Project Deploy Key Create Response JSON: %v", string(body))

	var created ProjectDeployKey
	if decodeerr := json.Unmarshal(body, &created); decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.SetId(fmt.Sprintf("%s/%s/%d", workspace, projectKey, created.ID))

	return resourceProjectDeployKeyRead(ctx, d, m)
}

func resourceProjectDeployKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace, projectKey, keyID, err := projectDeployKeyID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("2.0/workspaces/%s/projects/%s/deploy-keys/%s", workspace, projectKey, keyID)
	res, err := client.Get(url)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] Project Deploy Key (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	body, readerr := io.ReadAll(res.Body)
	if readerr != nil {
		return diag.FromErr(readerr)
	}

	var deployKey ProjectDeployKey
	if decodeerr := json.Unmarshal(body, &deployKey); decodeerr != nil {
		return diag.FromErr(decodeerr)
	}

	d.Set("workspace", workspace)
	d.Set("project_key", projectKey)
	d.Set("key_id", keyID)
	d.Set("label", deployKey.Label)
	d.Set("comment", deployKey.Comment)
	d.Set("added_on", deployKey.AddedOn)
	d.Set("last_used", deployKey.LastUsed)
	if deployKey.Key != "" {
		d.Set("key", deployKey.Key)
	}

	return nil
}

func resourceProjectDeployKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(Clients).httpClient

	workspace, projectKey, keyID, err := projectDeployKeyID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("2.0/workspaces/%s/projects/%s/deploy-keys/%s", workspace, projectKey, keyID)
	res, err := client.Delete(url)
	if err := handleClientError(res, err); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func projectDeployKeyID(id string) (string, string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%q), expected WORKSPACE/PROJECT-KEY/KEY-ID", id)
	}
	return parts[0], parts[1], parts[2], nil
}
