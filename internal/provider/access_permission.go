package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

var permissionsAttribute = &schema.Schema{
	Description: `The permissions of the access. Note: Not all permissions are available in all plans, and some permissions require others. More info in [1password docs](https://developer.1password.com/docs/cli/vault-permissions/).`,
	// TypeMap is currenty not supported in v2 sdk.
	Type:     schema.TypeList,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_viewing":           {Type: schema.TypeBool, Optional: true, Default: false},
			"allow_editing":           {Type: schema.TypeBool, Optional: true, Default: false},
			"allow_managing":          {Type: schema.TypeBool, Optional: true, Default: false},
			"view_items":              {Type: schema.TypeBool, Optional: true, Default: false},
			"create_items":            {Type: schema.TypeBool, Optional: true, Default: false},
			"edit_items":              {Type: schema.TypeBool, Optional: true, Default: false},
			"archive_items":           {Type: schema.TypeBool, Optional: true, Default: false},
			"delete_items":            {Type: schema.TypeBool, Optional: true, Default: false},
			"view_and_copy_passwords": {Type: schema.TypeBool, Optional: true, Default: false},
			"view_item_history":       {Type: schema.TypeBool, Optional: true, Default: false},
			"import_items":            {Type: schema.TypeBool, Optional: true, Default: false},
			"export_items":            {Type: schema.TypeBool, Optional: true, Default: false},
			"copy_and_share_items":    {Type: schema.TypeBool, Optional: true, Default: false},
			"print_items":             {Type: schema.TypeBool, Optional: true, Default: false},
			"manage_vault":            {Type: schema.TypeBool, Optional: true, Default: false},
		},
	},
}

func ensureDefaultValue(v interface{}) bool {
	if v == nil {
		return false
	}
	return v.(bool)
}

func dataToAccessPermissions(ap map[string]interface{}) model.AccessPermissions {
	return model.AccessPermissions{
		AllowViewing:         ensureDefaultValue(ap["allow_viewing"]),
		AllowEditing:         ensureDefaultValue(ap["allow_editing"]),
		AllowManaging:        ensureDefaultValue(ap["allow_managing"]),
		ViewItems:            ensureDefaultValue(ap["view_items"]),
		CreateItems:          ensureDefaultValue(ap["create_items"]),
		EditItems:            ensureDefaultValue(ap["edit_items"]),
		ArchiveItems:         ensureDefaultValue(ap["archive_items"]),
		DeleteItems:          ensureDefaultValue(ap["delete_items"]),
		ViewAndCopyPasswords: ensureDefaultValue(ap["view_and_copy_passwords"]),
		ViewItemHistory:      ensureDefaultValue(ap["view_item_history"]),
		ImportItems:          ensureDefaultValue(ap["import_items"]),
		ExportItems:          ensureDefaultValue(ap["export_items"]),
		CopyAndShareItems:    ensureDefaultValue(ap["copy_and_share_items"]),
		PrintItems:           ensureDefaultValue(ap["print_items"]),
		ManageVault:          ensureDefaultValue(ap["manage_vault"]),
	}
}

func accessPermissionsToData(m model.AccessPermissions) map[string]interface{} {
	return map[string]interface{}{
		"allow_viewing":           m.AllowViewing,
		"allow_editing":           m.AllowEditing,
		"allow_managing":          m.AllowManaging,
		"view_items":              m.ViewItems,
		"create_items":            m.CreateItems,
		"edit_items":              m.EditItems,
		"archive_items":           m.ArchiveItems,
		"delete_items":            m.DeleteItems,
		"view_and_copy_passwords": m.ViewAndCopyPasswords,
		"view_item_history":       m.ViewItemHistory,
		"import_items":            m.ImportItems,
		"export_items":            m.ExportItems,
		"copy_and_share_items":    m.CopyAndShareItems,
		"print_items":             m.PrintItems,
		"manage_vault":            m.ManageVault,
	}
}
