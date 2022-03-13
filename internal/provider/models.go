package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type User struct {
	ID    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
	Name  types.String `tfsdk:"name"`
}
