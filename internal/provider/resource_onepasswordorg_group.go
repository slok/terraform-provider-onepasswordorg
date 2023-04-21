package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: `
    Provides a Group resource.

    A 1password group is like a team that can contain people and can be used to give access to vaults as a
    group of users.
    `,
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,

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
				Description:  "The name of the group.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Managed by Terraform",
				Description: "The description of the group.",
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	group := dataToGroup(data)

	newGroup, err := p.repo.CreateGroup(ctx, group)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	groupToData(*newGroup, data)

	return diags
}

func resourceGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	// Get group.
	id := data.Id()
	group, err := p.repo.GetGroupByID(ctx, id)
	if err != nil {
		return diag.Errorf("Error reading group:" + fmt.Sprintf("Could not get group %q, unexpected error: %s", id, err.Error()))
	}

	groupToData(*group, data)
	return diags
}

func resourceGroupUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	id := data.Id()
	g := dataToGroup(data)

	newGroup, err := p.repo.EnsureGroup(ctx, g)
	if err != nil {
		return diag.Errorf("Error reading group:" + fmt.Sprintf("Could not get group %q, unexpected error: %s", id, err.Error()))
	}

	groupToData(*newGroup, data)
	return diags
}

func resourceGroupDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	// Get group.
	id := data.Id()
	err := p.repo.DeleteGroup(ctx, id)
	if err != nil {
		return diag.Errorf("Error reading group:" + fmt.Sprintf("Could not get group %q, unexpected error: %s", id, err.Error()))
	}

	return diags
}

func dataToGroup(data *schema.ResourceData) model.Group {
	return model.Group{
		ID:          data.Id(),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
	}
}

func groupToData(group model.Group, data *schema.ResourceData) {
	data.SetId(group.ID)
	data.Set("name", group.Name)
	data.Set("description", group.Description)
}
