package provider

import (
	"github.com/hashicorp/go-uuid"
	"strings"

	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func resourceItem() *schema.Resource {
	return &schema.Resource{
		Description:   "A 1Password item.",
		CreateContext: resourceItemCreate,
		ReadContext:   resourceItemRead,
		UpdateContext: resourceItemUpdate,
		DeleteContext: resourceItemDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The Terraform resource identifier for this item in the format `vaults/<vault_id>/items/<item_id>`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"uuid": {
				Description: itemUUIDDescription,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vault": {
				Description: vaultUUIDDescription,
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"category": {
				Description:  fmt.Sprintf(enumDescription, categoryDescription, categories),
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "login",
				ValidateFunc: validation.StringInSlice(categories, true),
				ForceNew:     true,
			},
			"title": {
				Description: itemTitleDescription,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"url": {
				Description: urlDescription,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"hostname": {
				Description: dbHostnameDescription,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"database": {
				Description: dbDatabaseDescription,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"port": {
				Description: dbPortDescription,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description:  fmt.Sprintf(enumDescription, dbTypeDescription, dbTypes),
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(dbTypes, true),
			},
			"tags": {
				Description: tagsDescription,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"username": {
				Description: usernameDescription,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password": {
				Description: passwordDescription,
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Computed:    true,
			},
			"section": {
				Description: sectionsDescription,
				Type:        schema.TypeList,
				Optional:    true,
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
							Required:    true,
						},
						"field": {
							Description: sectionFieldsDescription,
							Type:        schema.TypeList,
							Optional:    true,
							MinItems:    0,
							Elem: &schema.Resource{
								Description: fieldDescription,
								Schema: map[string]*schema.Schema{
									"id": {
										Description: fieldIDDescription,
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
									"label": {
										Description: fieldLabelDescription,
										Type:        schema.TypeString,
										Required:    true,
									},
									"purpose": {
										Description:  fmt.Sprintf(enumDescription, fieldPurposeDescription, fieldPurposes),
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(fieldPurposes, true),
									},
									"type": {
										Description:  fmt.Sprintf(enumDescription, fieldTypeDescription, fieldTypes),
										Type:         schema.TypeString,
										Default:      "STRING",
										Optional:     true,
										ValidateFunc: validation.StringInSlice(fieldTypes, true),
									},
									"value": {
										Description: fieldValueDescription,
										Type:        schema.TypeString,
										Optional:    true,
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

func resourceItemCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	item, err := dataToItem(data)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	newItem, err := p.repo.CreateItem(ctx, *item)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	itemToData(newItem, data)

	return diags
}

func resourceItemRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	// Get group.
	id := data.Id()
	_, itemUUID := vaultAndItemUUID(id)
	item, err := p.repo.GetItemByID(ctx, itemUUID)
	if err != nil {
		return diag.Errorf("Error reading group:" + fmt.Sprintf("Could not get item %q, unexpected error: %s", id, err.Error()))
	}

	itemToData(item, data)
	return diags
}

func resourceItemUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	id := data.Id()
	item, err := dataToItem(data)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	newItem, err := p.repo.EnsureItem(ctx, *item)
	if err != nil {
		return diag.Errorf("Error reading group:" + fmt.Sprintf("Could not get item %q, unexpected error: %s", id, err.Error()))
	}

	itemToData(newItem, data)
	return diags
}

func resourceItemDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	// Get group.
	id := data.Id()
	_, itemUUID := vaultAndItemUUID(id)
	err := p.repo.DeleteItem(ctx, itemUUID)
	if err != nil {
		return diag.Errorf("Error deleting item:" + fmt.Sprintf("Could not get item %q, unexpected error: %s", id, err.Error()))
	}

	return diags
}

func vaultAndItemUUID(tfID string) (vaultUUID, itemUUID string) {
	elements := strings.Split(tfID, "/")

	if len(elements) != 4 {
		return "", ""
	}

	return elements[1], elements[3]
}

func itemToData(item *model.Item, data *schema.ResourceData) {
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

	dataSections := data.Get("section").([]interface{})
	for _, s := range item.Sections {
		section := map[string]interface{}{}
		newSection := true

		// Check for existing section state
		for i := 0; i < len(dataSections); i++ {
			existingSection := dataSections[i].(map[string]interface{})
			existingID := existingSection["id"].(string)
			existingLabel := existingSection["label"].(string)

			if (s.ID != "" && s.ID == existingID) || s.Label == existingLabel {
				section = existingSection
				newSection = false
			}
		}

		section["id"] = s.ID
		section["label"] = s.Label

		existingFields := []interface{}{}
		if section["field"] != nil {
			existingFields = section["field"].([]interface{})
		}
		for _, f := range item.Fields {
			if f.Section != nil && f.Section.ID == s.ID {
				dataField := map[string]interface{}{}
				newField := true
				// Check for existing field state
				for i := 0; i < len(existingFields); i++ {
					existingField := existingFields[i].(map[string]interface{})
					existingID := existingField["id"].(string)
					existingLabel := existingField["label"].(string)

					if (f.ID != "" && f.ID == existingID) || f.Label == existingLabel {
						dataField = existingFields[i].(map[string]interface{})
						newField = false
					}
				}

				dataField["id"] = f.ID
				dataField["label"] = f.Label
				dataField["purpose"] = f.Purpose
				dataField["type"] = f.Type
				dataField["value"] = f.Value

				if newField {
					existingFields = append(existingFields, dataField)
				}
			}
		}
		section["field"] = existingFields

		if newSection {
			dataSections = append(dataSections, section)
		}
	}

	data.Set("section", dataSections)

	for _, f := range item.Fields {
		switch f.Purpose {
		case "USERNAME":
			data.Set("username", f.Value)
		case "PASSWORD":
			data.Set("password", f.Value)
		default:
			if f.Section == nil {
				data.Set(f.Label, f.Value)
			}
		}
	}
}

func dataToItem(data *schema.ResourceData) (*model.Item, error) {
	item := model.Item{
		ID: data.Get("uuid").(string),
		Vault: model.Vault{
			ID: data.Get("vault").(string),
		},
		Title: data.Get("title").(string),
		URLs: []model.URL{
			{
				Primary: true,
				URL:     data.Get("url").(string),
			},
		},
		Tags: getTags(data),
	}

	password := data.Get("password").(string)

	switch data.Get("category").(string) {
	case "login":
		item.Category = "login"
		item.Fields = []model.Field{
			{
				ID:      "username",
				Label:   "username",
				Purpose: "USERNAME",
				Type:    "STRING",
				Value:   data.Get("username").(string),
			},
			{
				ID:       "password",
				Label:    "password",
				Purpose:  "PASSWORD",
				Type:     "CONCEALED",
				Value:    password,
				Generate: password == "",
			},
		}
	case "password":
		item.Category = "password"
		item.Fields = []model.Field{
			{
				ID:       "password",
				Label:    "password",
				Purpose:  "PASSWORD",
				Type:     "CONCEALED",
				Value:    password,
				Generate: password == "",
			},
		}
	case "database":
		item.Category = "database"
		item.Fields = []model.Field{
			{
				ID:    "username",
				Label: "username",
				Type:  "STRING",
				Value: data.Get("username").(string),
			},
			{
				ID:       "password",
				Label:    "password",
				Type:     "CONCEALED",
				Value:    password,
				Generate: password == "",
			},
			{
				ID:    "hostname",
				Label: "hostname",
				Type:  "STRING",
				Value: data.Get("hostname").(string),
			},
			{
				ID:    "database",
				Label: "database",
				Type:  "STRING",
				Value: data.Get("database").(string),
			},
			{
				ID:    "port",
				Label: "port",
				Type:  "STRING",
				Value: data.Get("port").(string),
			},
			{
				ID:    "database_type",
				Label: "type",
				Type:  "MENU",
				Value: data.Get("type").(string),
			},
		}
	}

	sections := data.Get("section").([]interface{})
	for i := 0; i < len(sections); i++ {
		section, ok := sections[i].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unable to parse section: %v", sections[i])
		}
		sid, err := uuid.GenerateUUID()
		if err != nil {
			return nil, fmt.Errorf("Unable to generate a section id: %w", err)
		}

		if section["id"].(string) != "" {
			sid = section["id"].(string)
		} else {
			section["id"] = sid
		}

		s := &model.Section{
			ID:    sid,
			Label: section["label"].(string),
		}
		item.Sections = append(item.Sections, *s)

		sectionFields := section["field"].([]interface{})
		for j := 0; j < len(sectionFields); j++ {
			field, ok := sectionFields[j].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Unable to parse section field: %v", sectionFields[j])
			}

			f := &model.Field{
				Section: s,
				ID:      field["id"].(string),
				Type:    field["type"].(string),
				Purpose: field["purpose"].(string),
				Label:   field["label"].(string),
				Value:   field["value"].(string),
			}

			item.Fields = append(item.Fields, *f)
		}
	}

	return &item, nil
}

func getTags(data *schema.ResourceData) []string {
	tagInterface := data.Get("tags").([]interface{})
	tags := make([]string, len(tagInterface))
	for i, tag := range tagInterface {
		tags[i] = tag.(string)
	}
	return tags
}
