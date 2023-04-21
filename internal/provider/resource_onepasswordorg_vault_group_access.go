package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func resourceVaultGroupAccess() *schema.Resource {
	return &schema.Resource{
		Description: `
Provides vault access for a group.
    `,
		CreateContext: resourceVaultGroupAccessCreate,
		ReadContext:   resourceVaultGroupAccessRead,
		UpdateContext: resourceVaultGroupAccessUpdate,
		DeleteContext: resourceVaultGroupAccessDelete,

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
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The group ID.",
			},
			"permissions": permissionsAttribute,
		},
	}
}

func resourceVaultGroupAccessCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	m, err := dataToVaultGroupAccess(data)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	err = p.repo.EnsureVaultGroupAccess(ctx, *m)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	vaultGroupAccessToData(*m, data)

	return diags
}

func resourceVaultGroupAccessRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	id := data.Id()
	vaultID, groupID, err := unpackVaultGroupAccessID(id)
	if err != nil {
		return diag.Errorf("Error getting member ID: " + "Could not get member ID:" + err.Error())
	}

	vaultGroupAccess, err := p.repo.GetVaultGroupAccessByID(ctx, vaultID, groupID)
	if err != nil {
		return diag.Errorf("Error reading group access:" + fmt.Sprintf("Could not get group access %q, unexpected error: %s", id, err.Error()))
	}

	vaultGroupAccessToData(*vaultGroupAccess, data)
	return diags
}

func resourceVaultGroupAccessUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	m, err := dataToVaultGroupAccess(data)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	err = p.repo.EnsureVaultGroupAccess(ctx, *m)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	vaultGroupAccessToData(*m, data)

	return diags
}

func resourceVaultGroupAccessDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	// Get group access.
	id := data.Id()
	vaultID, groupID, err := unpackVaultGroupAccessID(id)
	if err != nil {
		return diag.Errorf("Error getting member ID: " + "Could not get member ID:" + err.Error())
	}

	err = p.repo.DeleteVaultGroupAccess(ctx, vaultID, groupID)
	if err != nil {
		return diag.Errorf("Error reading group access:" + fmt.Sprintf("Could not get group access %q, unexpected error: %s", id, err.Error()))
	}

	return diags
}

func dataToVaultGroupAccess(data *schema.ResourceData) (*model.VaultGroupAccess, error) {
	groupID := data.Get("group_id").(string)
	vaultID := data.Get("vault_id").(string)

	// Check the ID is correct.
	if data.Id() != "" {
		vid, gid, err := unpackVaultGroupAccessID(data.Id())
		if err != nil {
			return nil, err
		}

		if gid != groupID {
			return nil, fmt.Errorf("resource id is wrong based on group ID")
		}

		if vid != vaultID {
			return nil, fmt.Errorf("resource id is wrong based on vault ID")
		}
	}
	permissions := data.Get("permissions").([]interface{})

	return &model.VaultGroupAccess{
		VaultID:     vaultID,
		GroupID:     groupID,
		Permissions: dataToAccessPermissions(permissions[0].(map[string]interface{})),
	}, nil
}

func vaultGroupAccessToData(m model.VaultGroupAccess, data *schema.ResourceData) error {
	id := packVaultGroupAccessID(m.VaultID, m.GroupID)

	data.SetId(id)
	data.Set("group_id", m.GroupID)
	data.Set("vault_id", m.VaultID)
	data.Set("permissions", [1]map[string]interface{}{accessPermissionsToData(m.Permissions)})
	return nil
}

func packVaultGroupAccessID(vaultID, groupID string) string {
	return vaultID + "/" + groupID
}

func unpackVaultGroupAccessID(id string) (vaultID, groupID string, err error) {
	s := strings.SplitN(id, "/", 2)
	if len(s) != 2 {
		return "", "", fmt.Errorf(
			"invalid vault group access ID format: %s (expected <VAULT ID>/<GROUP ID>)", id)
	}

	return s[0], s[1], nil
}
