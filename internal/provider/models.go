package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type User struct {
	ID    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
	Name  types.String `tfsdk:"name"`
}

type Item struct {
	ID      types.String `tfsdk:"id"`
	VaultID types.String `tfsdk:"vault"`
	Title   types.String `tfsdk:"title"`
	Section []Section    `tfsdk:"section"`
}

type Section struct {
	ID    types.String `tfsdk:"id"`
	Field []Field      `tfsdk:"field"`
}

type Field struct {
	ID      types.String `tfsdk:"id"`
	Label   types.String `tfsdk:"label"`
	Type    types.String `tfsdk:"type"`
	Value   types.String `tfsdk:"value"`
	Purpose types.String `tfsdk:"purpose"`
}

type Group struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type Vault struct {
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
