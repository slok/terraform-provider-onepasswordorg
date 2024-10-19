package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

func New() provider.Provider {
	return &onePasswordOrgProvider{}
}

type onePasswordOrgProvider struct{}

func (p *onePasswordOrgProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "onepasswordorg"
}

func (p *onePasswordOrgProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
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
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Optional:    true,
				Description: fmt.Sprintf("Set account 1password domain address (e.g: something.1password.com). Also `%s` env var can be used.", envVarOpAddress),
			},
			"email": schema.StringAttribute{
				Optional:    true,
				Description: fmt.Sprintf("Set account 1password email. Also `%s` env var can be used.", envVarOpEmail),
			},
			"secret_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: fmt.Sprintf("Set account 1password secret key. Also `%s` env var can be used.", envVarOpSecretKey),
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: fmt.Sprintf("Set account 1password password. Also `%s` env var can be used.", envVarOpPassword),
			},
			"fake_storage_path": schema.StringAttribute{
				Optional:    true,
				Description: fmt.Sprintf("File to a path where the provider will store the data as if it is 1password (this is used only on development). Also `%s` env var can be used.", EnvVarOpFakeStoragePath),
			},
			"op_cli_path": schema.StringAttribute{
				Optional:    true,
				Description: fmt.Sprintf("The path that points to the op cli binary. Also `%s` env var can be used. (by default `op` on system path, ignored if run in Terraform cloud).", EnvVarOpCliPath),
			},
		},
	}
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

func (p *onePasswordOrgProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
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

	providerAppServices := providerAppServices{
		Repository: repo,
	}
	resp.DataSourceData = providerAppServices
	resp.ResourceData = providerAppServices
}

func (p *onePasswordOrgProvider) configureAddress(config providerData) (string, error) {
	if config.Address.IsUnknown() {
		return "", fmt.Errorf("cannot use unknown value as address")
	}

	// If not set get from env, the value has priority.
	var address string
	if config.Address.IsNull() {
		address = os.Getenv(envVarOpAddress)
	} else {
		address = config.Address.ValueString()
	}

	if address == "" {
		return "", fmt.Errorf("username cannot be an empty string")
	}

	return address, nil
}

func (p *onePasswordOrgProvider) configureEmail(config providerData) (string, error) {
	if config.Email.IsUnknown() {
		return "", fmt.Errorf("cannot use unknown value as email")
	}

	// If not set get from env, the value has priority.
	var email string
	if config.Email.IsNull() {
		email = os.Getenv(envVarOpEmail)
	} else {
		email = config.Email.ValueString()
	}

	if email == "" {
		return "", fmt.Errorf("email cannot be an empty string")
	}

	return email, nil
}

func (p *onePasswordOrgProvider) configureSecretKey(config providerData) (string, error) {
	if config.SecretKey.IsUnknown() {
		return "", fmt.Errorf("cannot use unknown value as secret key")
	}

	// If not set get from env, the value has priority.
	var secretKey string
	if config.SecretKey.IsNull() {
		secretKey = os.Getenv(envVarOpSecretKey)
	} else {
		secretKey = config.SecretKey.ValueString()
	}

	if secretKey == "" {
		return "", fmt.Errorf("secret key cannot be an empty string")
	}

	return secretKey, nil
}

func (p *onePasswordOrgProvider) configurePassword(config providerData) (string, error) {
	if config.Password.IsUnknown() {
		return "", fmt.Errorf("cannot use unknown value as password")
	}

	// If not set get from env, the value has priority.
	var password string
	if config.Password.IsNull() {
		password = os.Getenv(envVarOpPassword)
	} else {
		password = config.Password.ValueString()
	}

	if password == "" {
		return "", fmt.Errorf("password cannot be an empty string")
	}

	return password, nil
}

func (p *onePasswordOrgProvider) configureFakeStoragePath(config providerData) (string, error) {
	// If not set get from env, the value has priority.
	var fakePath string
	if config.FakeStoragePath.IsNull() {
		fakePath = os.Getenv(EnvVarOpFakeStoragePath)
	} else {
		fakePath = config.FakeStoragePath.ValueString()
	}

	return fakePath, nil
}

func (p *onePasswordOrgProvider) configureCliPath(config providerData) (string, error) {
	// If not set get from env, the value has priority.
	var cliPath string
	if config.FakeStoragePath.IsNull() {
		cliPath = os.Getenv(EnvVarOpCliPath)
	} else {
		cliPath = config.CliPath.ValueString()
	}

	return cliPath, nil
}

func (p *onePasswordOrgProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVaultResource,
		NewUserResource,
		NewGroupResource,
		NewGroupMemberResource,
		NewVaultUserAccessResource,
		NewVaultGroupAccessResource,
	}
}

func (p *onePasswordOrgProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVaultDataSource,
		NewUserDataSource,
		NewGroupDataSource,
	}
}

type providerAppServices struct {
	Repository storage.Repository
}
