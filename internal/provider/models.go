package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

type User struct {
	ID    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
	Name  types.String `tfsdk:"name"`
}

func mapModelToTfUser(u model.User) User {
	return User{
		ID:    types.StringValue(u.ID),
		Email: types.StringValue(u.Email),
		Name:  types.StringValue(u.Name),
	}
}

func mapTfToModelUser(u User) model.User {
	return model.User{
		ID:    u.ID.ValueString(),
		Email: u.Email.ValueString(),
		Name:  u.Name.ValueString(),
	}
}

type Group struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func mapModelToTfGroup(g model.Group) Group {
	return Group{
		ID:          types.StringValue(g.ID),
		Name:        types.StringValue(g.Name),
		Description: types.StringValue(g.Description),
	}
}

func mapTfToModelGroup(g Group) model.Group {
	return model.Group{
		ID:          g.ID.ValueString(),
		Name:        g.Name.ValueString(),
		Description: g.Description.ValueString(),
	}
}

type Vault struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func mapModelToTfVault(u model.Vault) Vault {
	return Vault{
		ID:          types.StringValue(u.ID),
		Name:        types.StringValue(u.Name),
		Description: types.StringValue(u.Description),
	}
}

func mapTfToModelVault(v Vault) model.Vault {
	return model.Vault{
		ID:          v.ID.ValueString(),
		Name:        v.Name.ValueString(),
		Description: v.Description.ValueString(),
	}
}

type Member struct {
	ID      types.String `tfsdk:"id"`
	UserID  types.String `tfsdk:"user_id"`
	GroupID types.String `tfsdk:"group_id"`
	Role    types.String `tfsdk:"role"`
}

type VaultGroupAccess struct {
	ID          types.String       `tfsdk:"id"`
	VaultID     types.String       `tfsdk:"vault_id"`
	GroupID     types.String       `tfsdk:"group_id"`
	Permissions *AccessPermissions `tfsdk:"permissions"`
}

type VaultUserAccess struct {
	ID          types.String       `tfsdk:"id"`
	VaultID     types.String       `tfsdk:"vault_id"`
	UserID      types.String       `tfsdk:"user_id"`
	Permissions *AccessPermissions `tfsdk:"permissions"`
}

type AccessPermissions struct {
	AllowViewing         types.Bool `tfsdk:"allow_viewing"`
	AllowEditing         types.Bool `tfsdk:"allow_editing"`
	AllowManaging        types.Bool `tfsdk:"allow_managing"`
	ViewItems            types.Bool `tfsdk:"view_items"`
	CreateItems          types.Bool `tfsdk:"create_items"`
	EditItems            types.Bool `tfsdk:"edit_items"`
	ArchiveItems         types.Bool `tfsdk:"archive_items"`
	DeleteItems          types.Bool `tfsdk:"delete_items"`
	ViewAndCopyPasswords types.Bool `tfsdk:"view_and_copy_passwords"`
	ViewItemHistory      types.Bool `tfsdk:"view_item_history"`
	ImportItems          types.Bool `tfsdk:"import_items"`
	ExportItems          types.Bool `tfsdk:"export_items"`
	CopyAndShareItems    types.Bool `tfsdk:"copy_and_share_items"`
	PrintItems           types.Bool `tfsdk:"print_items"`
	ManageVault          types.Bool `tfsdk:"manage_vault"`
}
