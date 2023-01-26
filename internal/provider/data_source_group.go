package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	group_name := data.Get("name").(string)
	group, err := p.repo.GetGroupByName(ctx, group_name)
	if err != nil {
		return diag.Errorf("Error getting group: Could not get group, unexpected error: " + err.Error())
	}

	data.SetId(group.ID)
	data.Set("name", group.Name)
	data.Set("description", group.Description)
	return diags
}

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: `
Provides information about a 1password group.
`,
		ReadContext: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the group.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"description": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"id": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}
