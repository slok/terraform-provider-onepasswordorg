package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: `
Provides information about a 1password user.
`,
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Description: "The email of the user.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"vaults": {
				Description: "List vaults that the user has access to.",
				Type:        schema.TypeList,
				Computed:    true,
				MinItems:    0,
				Elem: &schema.Resource{
					Description: sectionDescription,
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Id of the vault",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Name of the vault",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	email := data.Get("email").(string)
	user, err := p.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return diag.Errorf("Error getting user: Could not get user, unexpected error: " + err.Error())
	}

	vaults, err := p.repo.ListVaultsByUser(ctx, user.ID)
	if err != nil {
		return diag.Errorf("Error getting user: Could not get user, unexpected error: " + err.Error())
	}

	data.SetId(user.ID)
	data.Set("name", user.Name)
	data.Set("email", user.Email)

	dataVaults := []interface{}{}
	for _, s := range *vaults {
		vault := map[string]interface{}{}

		vault["id"] = s.ID
		vault["name"] = s.Name

		dataVaults = append(dataVaults, vault)
	}

	data.Set("vaults", dataVaults)
	return diags
}
