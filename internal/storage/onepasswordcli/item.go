package onepasswordcli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func (r Repository) CreateItem(ctx context.Context, item model.Item) (*model.Item, error) {
	cmdArgs := &onePasswordCliCmd{}
	// Create Fields
	cmdArgs.ItemArg().CreateArg()

	cmdArgs.EditFieldFlag("title", item.Title)
	cmdArgs.EditFieldFlag("url", item.URLs[0].URL)
	cmdArgs.CategoryFlag(item.Category)
	cmdArgs.VaultFlag(item.Vault.ID)

	for _, field := range item.Fields {
		if field.Section == nil {
			cmdArgs.RawStrArg(field.Label + "[" + field.Type + "]" + "=" + field.Value)
		} else {
			cmdArgs.RawStrArg(field.Section.Label + "." + field.Label + "[" + field.Type + "]" + "=" + field.Value)
		}
	}

	cmdArgs.FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ou := opItem{}
	err = json.Unmarshal([]byte(stdout), &ou)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotItem := mapOpToModelItem(ou)

	return &gotItem, nil
}

func (r Repository) GetItemByID(ctx context.Context, id string) (*model.Item, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.ItemArg().GetArg().RawStrArg(id).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ou := opItem{}
	err = json.Unmarshal([]byte(stdout), &ou)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotItem := mapOpToModelItem(ou)

	return &gotItem, nil
}
func (r Repository) GetItemByTitle(ctx context.Context, vaultID string, title string) (*model.Item, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.ItemArg().GetArg().RawStrArg(title).FormatJSONFlag().VaultFlag(vaultID)

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ou := opItem{}
	err = json.Unmarshal([]byte(stdout), &ou)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotItem := mapOpToModelItem(ou)

	return &gotItem, nil
}

func (r Repository) EnsureItem(ctx context.Context, item model.Item) (*model.Item, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.ItemArg().EditArg().RawStrArg(item.ID)

	cmdArgs.EditFieldFlag("title", item.Title)
	cmdArgs.EditFieldFlag("url", item.URLs[0].URL)
	cmdArgs.VaultFlag(item.Vault.ID)

	for _, field := range item.Fields {
		if field.Section == nil {
			cmdArgs.RawStrArg(field.Label + "[" + field.Type + "]" + "=" + field.Value)
		} else {
			cmdArgs.RawStrArg(field.Section.Label + "." + field.Label + "[" + field.Type + "]" + "=" + field.Value)
		}
	}

	cmdArgs.FormatJSONFlag()

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return &item, nil
}

func vaultAndItemUUID(tfID string) (vaultUUID, itemUUID string) {
	elements := strings.Split(tfID, "/")

	if len(elements) != 4 {
		return "", ""
	}

	return elements[1], elements[3]
}

func (r Repository) DeleteItem(ctx context.Context, id string) error {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.ItemArg().DeleteArg().RawStrArg(id)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return nil
}

type opItemField struct {
	ID      string     `json:"id"`
	Type    string     `json:"type"`
	Purpose string     `json:"purpose"`
	Label   string     `json:"label"`
	Value   string     `json:"value"`
	Section *opSection `json:"section"`
}

type opSection struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type opItem struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Category string        `json:"category"`
	Vault    opVault       `json:"vault"`
	Fields   []opItemField `json:"fields"`
	Sections []opSection   `json:"sections"`
	Tags     []string      `json:"tags"`
}

func mapOpToModelItem(u opItem) model.Item {
	return model.Item{
		ID:       u.ID,
		Title:    u.Title,
		Tags:     u.Tags,
		Category: u.Category,
		Vault:    mapOpToModeVault(u.Vault),
		Fields:   mapOpToModelItemFields(u.Fields),
		Sections: mapOpToModelItemSections(u.Sections),
	}
}
func mapOpToModelSection(u *opSection) model.Section {
	return model.Section{
		ID:    u.ID,
		Label: u.Label,
	}
}

func mapOpToModelItemSections(opSections []opSection) []model.Section {
	sections := []model.Section{}

	for _, opSection := range opSections {
		section := model.Section{
			ID:    opSection.ID,
			Label: opSection.Label,
		}
		sections = append(sections, section)
	}

	return sections
}

func mapOpToModelItemFields(opFields []opItemField) []model.Field {
	fields := []model.Field{}

	for _, opField := range opFields {
		field := model.Field{
			ID:      opField.ID,
			Type:    opField.Type,
			Purpose: opField.Purpose,
			Label:   opField.Label,
			Value:   opField.Value,
		}
		if opField.Section != nil {
			section := mapOpToModelSection(opField.Section)
			field.Section = &section
		}
		fields = append(fields, field)
	}

	return fields
}
