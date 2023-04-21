package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func resourceVaultUserAccess() *schema.Resource {
	return &schema.Resource{
		Description: `
Provides vault access for a user.
    `,
		CreateContext: resourceVaultUserAccessCreate,
		ReadContext:   resourceVaultUserAccessRead,
		UpdateContext: resourceVaultUserAccessUpdate,
		DeleteContext: resourceVaultUserAccessDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vault_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The vault ID.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The user ID.",
			},
			"permissions": permissionsAttribute,
		},
	}
}

func resourceVaultUserAccessCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	m, err := dataToVaultUserAccess(data)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	err = p.repo.EnsureVaultUserAccess(ctx, *m)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	vaultUserAccessToData(*m, data)

	return diags
}

func resourceVaultUserAccessRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	id := data.Id()
	vaultID, userID, err := unpackVaultUserAccessID(id)
	if err != nil {
		return diag.Errorf("Error getting member ID: " + "Could not get member ID:" + err.Error())
	}

	vaultUserAccess, err := p.repo.GetVaultUserAccessByID(ctx, vaultID, userID)
	if err != nil {
		return diag.Errorf("Error reading user access:" + fmt.Sprintf("Could not get user access %q, unexpected error: %s", id, err.Error()))
	}

	vaultUserAccessToData(*vaultUserAccess, data)
	return diags
}

func resourceVaultUserAccessUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	m, err := dataToVaultUserAccess(data)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	err = p.repo.EnsureVaultUserAccess(ctx, *m)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	vaultUserAccessToData(*m, data)

	return diags
}

func resourceVaultUserAccessDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	// Get user access.
	id := data.Id()
	vaultID, userID, err := unpackVaultUserAccessID(id)
	if err != nil {
		return diag.Errorf("Error getting member ID: " + "Could not get member ID:" + err.Error())
	}

	err = p.repo.DeleteVaultUserAccess(ctx, vaultID, userID)
	if err != nil {
		return diag.Errorf("Error reading user access:" + fmt.Sprintf("Could not get user access %q, unexpected error: %s", id, err.Error()))
	}

	return diags
}

func dataToVaultUserAccess(data *schema.ResourceData) (*model.VaultUserAccess, error) {
	userID := data.Get("user_id").(string)
	vaultID := data.Get("vault_id").(string)

	// Check the ID is correct.
	if data.Id() != "" {
		vid, gid, err := unpackVaultUserAccessID(data.Id())
		if err != nil {
			return nil, err
		}

		if gid != userID {
			return nil, fmt.Errorf("resource id is wrong based on user ID")
		}

		if vid != vaultID {
			return nil, fmt.Errorf("resource id is wrong based on vault ID")
		}
	}
	permissions := data.Get("permissions").([]interface{})

	return &model.VaultUserAccess{
		VaultID:     vaultID,
		UserID:      userID,
		Permissions: dataToAccessPermissions(permissions[0].(map[string]interface{})),
	}, nil
}

func vaultUserAccessToData(m model.VaultUserAccess, data *schema.ResourceData) error {
	id := packVaultUserAccessID(m.VaultID, m.UserID)

	data.SetId(id)
	data.Set("user_id", m.UserID)
	data.Set("vault_id", m.VaultID)
	data.Set("permissions", [1]map[string]interface{}{accessPermissionsToData(m.Permissions)})
	return nil
}
func packVaultUserAccessID(vaultID, userID string) string {
	return vaultID + "/" + userID
}

func unpackVaultUserAccessID(id string) (vaultID, userID string, err error) {
	s := strings.SplitN(id, "/", 2)
	if len(s) != 2 {
		return "", "", fmt.Errorf(
			"invalid vault user access ID format: %s (expected <VAULT ID>/<USER ID>)", id)
	}

	return s[0], s[1], nil
}
