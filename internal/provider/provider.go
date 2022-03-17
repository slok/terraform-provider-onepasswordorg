package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/slok/terraform-provider-onepasswordorg/internal/storage"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage/fake"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage/onepasswordcli"
)

const (
	envVarOpAddress         = "OP_ADDRESS"
	envVarOpEmail           = "OP_EMAIL"
	envVarOpSecretKey       = "OP_SECRET_KEY"
	envVarOpPassword        = "OP_PASSWORD"
	EnvVarOpFakeStoragePath = "OP_FAKE_STORAGE_PATH"
	EnvVarOpCliPath         = "OP_CLI_PATH"
)

func New() tfsdk.Provider {
	return &provider{}

}

type provider struct {
	configured bool
	repo       storage.Repository
}

// GetSchema returns the schema that the user must configure on the provider block.
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
The Onepassword organization provider is used to interact with 1password organization resources (users, groups...)
and not items.

Normally this provider will be used to automate the user and groups management like user onboard/offboards or
grouping users into teams (groups in 1password).

## Requirements

This provider needs [op](https://1password.com/downloads/command-line/) v2.x Cli, thats why it doesn't use 1password connect
API and needs a real 1password account as the authentication.

## Authentication

Needs a real 1password account so the provider can use the "password" and "secret key" of that account.

A recommended way would be creating an account in the 1password organization/company only for automation
like Terraform (used by this provider).

## Terraform cloud

The provider will detect that its executing in terraform cloud and will use the embedded op CLI for this purpose
so it satisfies the op Cli requirement inside Terraform cloud workers.
`,
		Attributes: map[string]tfsdk.Attribute{
			"address": {
				Type:        types.StringType,
				Optional:    true,
				Description: fmt.Sprintf("Set account 1password domain address (e.g: something.1password.com). Also `%s` env var can be used.", envVarOpAddress),
			},
			"email": {
				Type:        types.StringType,
				Optional:    true,
				Description: fmt.Sprintf("Set account 1password email. Also `%s` env var can be used.", envVarOpEmail),
			},
			"secret_key": {
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
				Description: fmt.Sprintf("Set account 1password secret key. Also `%s` env var can be used.", envVarOpSecretKey),
			},
			"password": {
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
				Description: fmt.Sprintf("Set account 1password password. Also `%s` env var can be used.", envVarOpPassword),
			},
			"fake_storage_path": {
				Type:        types.StringType,
				Optional:    true,
				Description: fmt.Sprintf("File to a path where the provider will store the data as if it is 1password (this is used only on development). Also `%s` env var can be used.", EnvVarOpFakeStoragePath),
			},
			"op_cli_path": {
				Type:        types.StringType,
				Optional:    true,
				Description: fmt.Sprintf("The path that points to the op cli binary. Also `%s` env var can be used. (by default `op` on system path, ignored if run in Terraform cloud).", EnvVarOpCliPath),
			},
		},
	}, nil
}

// Provider configuration.
type providerData struct {
	Address         types.String `tfsdk:"address"`
	Email           types.String `tfsdk:"email"`
	SecretKey       types.String `tfsdk:"secret_key"`
	Password        types.String `tfsdk:"password"`
	FakeStoragePath types.String `tfsdk:"fake_storage_path"`
	CliPath         types.String `tfsdk:"op_cli_path"`
}

// This is like if it was our main entrypoint.
func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration.
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Error summaries
	const (
		configErrSummary = "Unable to configure client"
		createErrSummary = "Unable to create op client"
	)

	// Get if we are in fake mode.
	fakeStoragePath, err := p.configureFakeStoragePath(config)
	if err != nil {
		resp.Diagnostics.AddError(configErrSummary, "Invalid fake storage path:\n\n"+err.Error())
	}

	// Create fake or regular mode.
	// If the user has set the fake storage path then we are going to use a fake repository.
	// If the user didn't, we will use the op cli based repository (a.k.a real 1password APIs).
	var repo storage.Repository
	if fakeStoragePath != "" {
		repo, err = fake.NewRepository(fakeStoragePath)
		if err != nil {
			resp.Diagnostics.AddError(createErrSummary, "Unable to create 1password fake storage:\n\n"+err.Error())
			return
		}
	} else {
		address, err := p.configureAddress(config)
		if err != nil {
			resp.Diagnostics.AddError(configErrSummary, "Invalid address:\n\n"+err.Error())
		}

		email, err := p.configureEmail(config)
		if err != nil {
			resp.Diagnostics.AddError(configErrSummary, "Invalid email:\n\n"+err.Error())
		}

		secretKey, err := p.configureSecretKey(config)
		if err != nil {
			resp.Diagnostics.AddError(configErrSummary, "Invalid secret key:\n\n"+err.Error())
		}

		password, err := p.configurePassword(config)
		if err != nil {
			resp.Diagnostics.AddError(configErrSummary, "Invalid password:\n\n"+err.Error())
		}

		cliPath, err := p.configureCliPath(config)
		if err != nil {
			resp.Diagnostics.AddError(configErrSummary, "Invalid cli path:\n\n"+err.Error())
		}

		// Create OP cli.
		cli, err := onepasswordcli.NewOpCli(cliPath, address, email, secretKey, password)
		if err != nil {
			resp.Diagnostics.AddError(createErrSummary, "Unable to create 1password op cmd client:\n\n"+err.Error())
			return
		}

		// Create  repository.
		repo, err = onepasswordcli.NewRepository(cli)
		if err != nil {
			resp.Diagnostics.AddError(createErrSummary, "Unable to create 1password op repository:\n\n"+err.Error())
			return
		}
	}

	p.repo = repo
	p.configured = true
}

func (p *provider) configureAddress(config providerData) (string, error) {
	if config.Address.Unknown {
		return "", fmt.Errorf("cannot use unknown value as address")
	}

	// If not set get from env, the value has priority.
	var address string
	if config.Address.Null {
		address = os.Getenv(envVarOpAddress)
	} else {
		address = config.Address.Value
	}

	if address == "" {
		return "", fmt.Errorf("username cannot be an empty string")
	}

	return address, nil
}

func (p *provider) configureEmail(config providerData) (string, error) {
	if config.Email.Unknown {
		return "", fmt.Errorf("cannot use unknown value as email")
	}

	// If not set get from env, the value has priority.
	var email string
	if config.Email.Null {
		email = os.Getenv(envVarOpEmail)
	} else {
		email = config.Email.Value
	}

	if email == "" {
		return "", fmt.Errorf("email cannot be an empty string")
	}

	return email, nil
}

func (p *provider) configureSecretKey(config providerData) (string, error) {
	if config.SecretKey.Unknown {
		return "", fmt.Errorf("cannot use unknown value as secret key")
	}

	// If not set get from env, the value has priority.
	var secretKey string
	if config.SecretKey.Null {
		secretKey = os.Getenv(envVarOpSecretKey)
	} else {
		secretKey = config.SecretKey.Value
	}

	if secretKey == "" {
		return "", fmt.Errorf("secret key cannot be an empty string")
	}

	return secretKey, nil
}

func (p *provider) configurePassword(config providerData) (string, error) {
	if config.Password.Unknown {
		return "", fmt.Errorf("cannot use unknown value as password")
	}

	// If not set get from env, the value has priority.
	var password string
	if config.Password.Null {
		password = os.Getenv(envVarOpPassword)
	} else {
		password = config.Password.Value
	}

	if password == "" {
		return "", fmt.Errorf("password cannot be an empty string")
	}

	return password, nil
}

func (p *provider) configureFakeStoragePath(config providerData) (string, error) {
	// If not set get from env, the value has priority.
	var fakePath string
	if config.FakeStoragePath.Null {
		fakePath = os.Getenv(EnvVarOpFakeStoragePath)
	} else {
		fakePath = config.FakeStoragePath.Value
	}

	return fakePath, nil
}

func (p *provider) configureCliPath(config providerData) (string, error) {
	// If not set get from env, the value has priority.
	var cliPath string
	if config.FakeStoragePath.Null {
		cliPath = os.Getenv(EnvVarOpCliPath)
	} else {
		cliPath = config.CliPath.Value
	}

	return cliPath, nil
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"onepasswordorg_user":               resourceUserType{},
		"onepasswordorg_group":              resourceGroupType{},
		"onepasswordorg_group_member":       resourceGroupMemberType{},
		"onepasswordorg_vault":              resourceVaultType{},
		"onepasswordorg_vault_group_access": resourceVaultGroupAccessType{},
	}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"onepasswordorg_user":  dataSourceUserType{},
		"onepasswordorg_group": dataSourceGroupType{},
		"onepasswordorg_vault": dataSourceVaultType{},
	}, nil
}
