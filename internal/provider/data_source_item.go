package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

const (
	itemUUIDDescription  = "The UUID of the item. Item identifiers are unique within a specific vault."
	vaultUUIDDescription = "The UUID of the vault the item is in."
	categoryDescription  = "The category of the item."
	itemTitleDescription = "The title of the item."
	urlDescription       = "The primary URL for the item."
	tagsDescription      = "An array of strings of the tags assigned to the item."
	usernameDescription  = "Username for this item."
	passwordDescription  = "Password for this item."
	noteValueDescription = "Secure Note value."

	dbHostnameDescription = "(Only applies to the database category) The address where the database can be found"
	dbDatabaseDescription = "(Only applies to the database category) The name of the database."
	dbPortDescription     = "(Only applies to the database category) The port the database is listening on."
	dbTypeDescription     = "(Only applies to the database category) The type of database."

	sectionsDescription      = "A list of custom sections in an item"
	sectionDescription       = "A custom section in an item that contains custom fields"
	sectionIDDescription     = "A unique identifier for the section."
	sectionLabelDescription  = "The label for the section."
	sectionFieldsDescription = "A list of custom fields in the section."

	fieldDescription        = "A custom field."
	fieldIDDescription      = "A unique identifier for the field."
	fieldLabelDescription   = "The label for the field."
	fieldPurposeDescription = "Purpose indicates this is a special field: a username, password, or notes field."
	fieldTypeDescription    = "The type of value stored in the field."
	fieldValueDescription   = "The value of the field."

	passwordRecipeDescription  = "The recipe used to generate a new value for a password."
	passwordElementDescription = "The kinds of characters to include in the password."
	passwordLengthDescription  = "The length of the password to be generated."
	passwordLettersDescription = "Use letters [a-zA-Z] when generating the password."
	passwordDigitsDescription  = "Use digits [0-9] when generating the password."
	passwordSymbolsDescription = "Use symbols [!@.-_*] when generating the password."

	enumDescription = "%s One of %q"
)

var categories = []string{"login", "password", "database"}
var dbTypes = []string{"db2", "filemaker", "msaccess", "mssql", "mysql", "oracle", "postgresql", "sqlite", "other"}
var fieldPurposes = []string{"USERNAME", "PASSWORD", "NOTES"}
var fieldTypes = []string{"STRING", "EMAIL", "CONCEALED", "URL", "OTP", "DATE", "MONTH_YEAR", "MENU"}

func dataSourceItem() *schema.Resource {
	exactlyOneOfUUIDAndTitle := []string{"uuid", "title"}

	return &schema.Resource{
		Description: "Use this data source to get details of an item by its vault uuid and either the title or the uuid of the item.",
		ReadContext: dataSourceItemRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The Terraform resource identifier for this item in the format `vaults/<vault_id>/items/<item_id>`",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vault": {
				Description: vaultUUIDDescription,
				Type:        schema.TypeString,
				Required:    true,
			},
			"uuid": {
				Description:  "The UUID of the item to retrieve. This field will be populated with the UUID of the item if the item it looked up by its title.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: exactlyOneOfUUIDAndTitle,
			},
			"title": {
				Description:  "The title of the item to retrieve. This field will be populated with the title of the item if the item it looked up by its UUID.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: exactlyOneOfUUIDAndTitle,
			},
			"category": {
				Description: fmt.Sprintf(enumDescription, categoryDescription, categories),
				Type:        schema.TypeString,
				Computed:    true,
			},
			"url": {
				Description: urlDescription,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"hostname": {
				Description: dbHostnameDescription,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"database": {
				Description: dbDatabaseDescription,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"port": {
				Description: dbPortDescription,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: fmt.Sprintf(enumDescription, dbTypeDescription, dbTypes),
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tags": {
				Description: tagsDescription,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"username": {
				Description: usernameDescription,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"password": {
				Description: passwordDescription,
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"note_value": {
				Description: noteValueDescription,
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Sensitive:   true,
			},
			"section": {
				Description: sectionsDescription,
				Type:        schema.TypeList,
				Computed:    true,
				MinItems:    0,
				Elem: &schema.Resource{
					Description: sectionDescription,
					Schema: map[string]*schema.Schema{
						"id": {
							Description: sectionIDDescription,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"label": {
							Description: sectionLabelDescription,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"field": {
							Description: sectionFieldsDescription,
							Type:        schema.TypeList,
							Computed:    true,
							MinItems:    0,
							Elem: &schema.Resource{
								Description: fieldDescription,
								Schema: map[string]*schema.Schema{
									"id": {
										Description: fieldIDDescription,
										Type:        schema.TypeString,
										Computed:    true,
									},
									"label": {
										Description: fieldLabelDescription,
										Type:        schema.TypeString,
										Computed:    true,
									},
									"purpose": {
										Description: fmt.Sprintf(enumDescription, fieldPurposeDescription, fieldPurposes),
										Type:        schema.TypeString,
										Computed:    true,
									},
									"type": {
										Description: fmt.Sprintf(enumDescription, fieldTypeDescription, fieldTypes),
										Type:        schema.TypeString,
										Computed:    true,
									},
									"value": {
										Description: fieldValueDescription,
										Type:        schema.TypeString,
										Computed:    true,
										Sensitive:   true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func terraformID(item model.Item) string {
	return fmt.Sprintf("vaults/%s/items/%s", item.Vault.ID, item.ID)
}

func dataSourceItemRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)

	var diags diag.Diagnostics

	vaultUUID := data.Get("vault").(string)
	itemTitle := data.Get("title").(string)
	itemUUID := data.Get("uuid").(string)

	query := ""
	if itemTitle != "" {
		query = itemTitle
	}
	if itemUUID != "" {
		query = itemUUID
	}

	item, err := p.repo.GetItemByTitle(ctx, vaultUUID, query)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	data.SetId(terraformID(*item))
	data.Set("uuid", item.ID)
	data.Set("vault", item.Vault.ID)
	data.Set("title", item.Title)

	for _, u := range item.URLs {
		if u.Primary {
			data.Set("url", u.URL)
		}
	}

	data.Set("tags", item.Tags)
	data.Set("category", strings.ToLower(string(item.Category)))

	dataSections := []interface{}{}
	for _, s := range item.Sections {
		section := map[string]interface{}{}

		section["id"] = s.ID
		section["label"] = s.Label

		fields := []interface{}{}

		for _, f := range item.Fields {
			if f.Section != nil && f.Section.ID == s.ID {
				dataField := map[string]interface{}{}
				dataField["id"] = f.ID
				dataField["label"] = f.Label
				dataField["purpose"] = f.Purpose
				dataField["type"] = f.Type
				dataField["value"] = f.Value

				fields = append(fields, dataField)
			}
		}
		section["field"] = fields

		dataSections = append(dataSections, section)
	}

	data.Set("section", dataSections)

	for _, f := range item.Fields {
		switch f.Purpose {
		case "USERNAME":
			data.Set("username", f.Value)
		case "PASSWORD":
			data.Set("password", f.Value)
		case "NOTES":
			data.Set("note_value", f.Value)
		default:
			if f.Section == nil {
				data.Set(strings.ToLower(f.Label), f.Value)
			}
		}
	}

	return diags
}
