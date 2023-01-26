package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description: `
Provides a User resource.

When a 1password user resources is created, it will be invited  by email.
    `,
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

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
				Description:  "The name of the user.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The description of the user.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	user := dataToUser(data)

	newUser, err := p.repo.CreateUser(ctx, user)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	userToData(*newUser, data)

	return diags
}

func resourceUserRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	// Get user.
	id := data.Id()
	user, err := p.repo.GetUserByID(ctx, id)
	if err != nil {
		return diag.Errorf("Error reading user:" + fmt.Sprintf("Could not get user %q, unexpected error: %s", id, err.Error()))
	}

	userToData(*user, data)
	return diags
}

func resourceUserUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	id := data.Id()
	g := dataToUser(data)

	newUser, err := p.repo.EnsureUser(ctx, g)
	if err != nil {
		return diag.Errorf("Error reading user:" + fmt.Sprintf("Could not get user %q, unexpected error: %s", id, err.Error()))
	}

	userToData(*newUser, data)
	return diags
}

func resourceUserDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	// Get user.
	id := data.Id()
	err := p.repo.DeleteUser(ctx, id)
	if err != nil {
		return diag.Errorf("Error reading user:" + fmt.Sprintf("Could not get user %q, unexpected error: %s", id, err.Error()))
	}

	return diags
}

func dataToUser(data *schema.ResourceData) model.User {
	return model.User{
		ID:    data.Id(),
		Name:  data.Get("name").(string),
		Email: data.Get("email").(string),
	}
}

func userToData(user model.User, data *schema.ResourceData) {
	data.SetId(user.ID)
	data.Set("name", user.Name)
	data.Set("email", user.Email)
}
