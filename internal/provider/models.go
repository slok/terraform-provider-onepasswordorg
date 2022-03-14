package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type User struct {
	ID    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
	Name  types.String `tfsdk:"name"`
}

type Group struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type Member struct {
	ID      types.String `tfsdk:"id"`
	UserID  types.String `tfsdk:"user_id"`
	GroupID types.String `tfsdk:"group_id"`
	Role    types.String `tfsdk:"role"`
}
