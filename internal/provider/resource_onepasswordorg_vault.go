package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

type resourceVaultType struct{}

func resourceVault() *schema.Resource {
	return &schema.Resource{
		Description: `
Provides a vault resource.
    `,
		CreateContext: resourceVaultCreate,
		ReadContext:   resourceVaultRead,
		UpdateContext: resourceVaultUpdate,
		DeleteContext: resourceVaultDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the vault.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Managed by Terraform",
				Description: "The description of the vault.",
			},
		},
	}
}

func resourceVaultCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	vault := dataToVault(data)

	newVault, err := p.repo.CreateVault(ctx, vault)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	vaultToData(*newVault, data)

	return diags
}

func resourceVaultRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	// Get vault.
	id := data.Id()
	vault, err := p.repo.GetVaultByID(ctx, id)
	if err != nil {
		return diag.Errorf("Error reading vault:" + fmt.Sprintf("Could not get vault %q, unexpected error: %s", id, err.Error()))
	}

	vaultToData(*vault, data)
	return diags
}

func resourceVaultUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	id := data.Id()
	vault := dataToVault(data)

	newVault, err := p.repo.EnsureVault(ctx, vault)
	if err != nil {
		return diag.Errorf("Error reading vault:" + fmt.Sprintf("Could not get vault %q, unexpected error: %s", id, err.Error()))
	}

	vaultToData(*newVault, data)
	return diags
}

func resourceVaultDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	// Get vault.
	id := data.Id()
	err := p.repo.DeleteVault(ctx, id)
	if err != nil {
		return diag.Errorf("Error reading vault:" + fmt.Sprintf("Could not get vault %q, unexpected error: %s", id, err.Error()))
	}

	return diags
}

func dataToVault(data *schema.ResourceData) model.Vault {
	return model.Vault{
		ID:          data.Id(),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
	}
}

func vaultToData(vault model.Vault, data *schema.ResourceData) {
	data.SetId(vault.ID)
	data.Set("name", vault.Name)
	data.Set("description", vault.Description)
}
