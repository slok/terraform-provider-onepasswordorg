package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

var permissionsAttribute = schema.SingleNestedAttribute{
	Required:    true,
	Description: `The permissions of the access. Note: Not all permissions are available in all plans, and some permissions require others. More info in [1password docs](https://developer.1password.com/docs/cli/vault-permissions/).`,
	Attributes: map[string]schema.Attribute{
		"allow_viewing":           schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"allow_editing":           schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"allow_managing":          schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"view_items":              schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"create_items":            schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"edit_items":              schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"archive_items":           schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"delete_items":            schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"view_and_copy_passwords": schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"view_item_history":       schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"import_items":            schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"export_items":            schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"copy_and_share_items":    schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"print_items":             schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
		"manage_vault":            schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
	},
}

func mapTfToModelAccessPermissions(ap AccessPermissions) model.AccessPermissions {
	return model.AccessPermissions{
		AllowViewing:         ap.AllowViewing.ValueBool(),
		AllowEditing:         ap.AllowEditing.ValueBool(),
		AllowManaging:        ap.AllowManaging.ValueBool(),
		ViewItems:            ap.ViewItems.ValueBool(),
		CreateItems:          ap.CreateItems.ValueBool(),
		EditItems:            ap.EditItems.ValueBool(),
		ArchiveItems:         ap.ArchiveItems.ValueBool(),
		DeleteItems:          ap.DeleteItems.ValueBool(),
		ViewAndCopyPasswords: ap.ViewAndCopyPasswords.ValueBool(),
		ViewItemHistory:      ap.ViewItemHistory.ValueBool(),
		ImportItems:          ap.ImportItems.ValueBool(),
		ExportItems:          ap.ExportItems.ValueBool(),
		CopyAndShareItems:    ap.CopyAndShareItems.ValueBool(),
		PrintItems:           ap.PrintItems.ValueBool(),
		ManageVault:          ap.ManageVault.ValueBool(),
	}
}

func mapModelToTfAccessPermissions(m model.AccessPermissions) *AccessPermissions {
	return &AccessPermissions{
		AllowViewing:         types.BoolValue(m.AllowViewing),
		AllowEditing:         types.BoolValue(m.AllowEditing),
		AllowManaging:        types.BoolValue(m.AllowManaging),
		ViewItems:            types.BoolValue(m.ViewItems),
		CreateItems:          types.BoolValue(m.CreateItems),
		EditItems:            types.BoolValue(m.EditItems),
		ArchiveItems:         types.BoolValue(m.ArchiveItems),
		DeleteItems:          types.BoolValue(m.DeleteItems),
		ViewAndCopyPasswords: types.BoolValue(m.ViewAndCopyPasswords),
		ViewItemHistory:      types.BoolValue(m.ViewItemHistory),
		ImportItems:          types.BoolValue(m.ImportItems),
		ExportItems:          types.BoolValue(m.ExportItems),
		CopyAndShareItems:    types.BoolValue(m.CopyAndShareItems),
		PrintItems:           types.BoolValue(m.PrintItems),
		ManageVault:          types.BoolValue(m.ManageVault),
	}
}
