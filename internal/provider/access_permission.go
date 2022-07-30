package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider/attributeutils"
)

var permissionsAttribute = tfsdk.Attribute{
	Required:    true,
	Description: `The permissions of the access. Note: Not all permissions are available in all plans, and some permissions require others. More info in [1password docs](https://developer.1password.com/docs/cli/vault-permissions/).`,
	Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
		"allow_viewing":           {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"allow_editing":           {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"allow_managing":          {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"view_items":              {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"create_items":            {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"edit_items":              {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"archive_items":           {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"delete_items":            {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"view_and_copy_passwords": {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"view_item_history":       {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"import_items":            {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"export_items":            {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"copy_and_share_items":    {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"print_items":             {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
		"manage_vault":            {Type: types.BoolType, Computed: true, Optional: true, PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})}},
	}),
}

func mapTfToModelAccessPermissions(ap AccessPermissions) model.AccessPermissions {
	return model.AccessPermissions{
		AllowViewing:         ap.AllowViewing.Value,
		AllowEditing:         ap.AllowEditing.Value,
		AllowManaging:        ap.AllowManaging.Value,
		ViewItems:            ap.ViewItems.Value,
		CreateItems:          ap.CreateItems.Value,
		EditItems:            ap.EditItems.Value,
		ArchiveItems:         ap.ArchiveItems.Value,
		DeleteItems:          ap.DeleteItems.Value,
		ViewAndCopyPasswords: ap.ViewAndCopyPasswords.Value,
		ViewItemHistory:      ap.ViewItemHistory.Value,
		ImportItems:          ap.ImportItems.Value,
		ExportItems:          ap.ExportItems.Value,
		CopyAndShareItems:    ap.CopyAndShareItems.Value,
		PrintItems:           ap.PrintItems.Value,
		ManageVault:          ap.ManageVault.Value,
	}
}

func mapModelToTfAccessPermissions(m model.AccessPermissions) *AccessPermissions {
	return &AccessPermissions{
		AllowViewing:         types.Bool{Value: m.AllowViewing},
		AllowEditing:         types.Bool{Value: m.AllowEditing},
		AllowManaging:        types.Bool{Value: m.AllowManaging},
		ViewItems:            types.Bool{Value: m.ViewItems},
		CreateItems:          types.Bool{Value: m.CreateItems},
		EditItems:            types.Bool{Value: m.EditItems},
		ArchiveItems:         types.Bool{Value: m.ArchiveItems},
		DeleteItems:          types.Bool{Value: m.DeleteItems},
		ViewAndCopyPasswords: types.Bool{Value: m.ViewAndCopyPasswords},
		ViewItemHistory:      types.Bool{Value: m.ViewItemHistory},
		ImportItems:          types.Bool{Value: m.ImportItems},
		ExportItems:          types.Bool{Value: m.ExportItems},
		CopyAndShareItems:    types.Bool{Value: m.CopyAndShareItems},
		PrintItems:           types.Bool{Value: m.PrintItems},
		ManageVault:          types.Bool{Value: m.ManageVault},
	}
}
