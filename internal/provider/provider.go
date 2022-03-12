package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/slok/terraform-provider-onepasswordorg/internal/storage/onepasswordcli"
)

const (
	envVarOpAddress   = "OP_ADDRESS"
	envVarOpEmail     = "OP_EMAIL"
	envVarOpSecretKey = "OP_SECRET_KEY"
)

var stderr = os.Stderr

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	repo       *onepasswordcli.Repository
}

// GetSchema returns the schema that the user must configure on the provider block.
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"address": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Set account 1password domain address (e.g: something.1password.com)",
			},
			"email": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Set account 1password email",
			},
			"secret_key": {
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
				Description: "Set account 1password secret key",
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct {
	Address   types.String `tfsdk:"address"`
	Email     types.String `tfsdk:"email"`
	SecretKey types.String `tfsdk:"secret_key"`
}

// This is like if it was our main entrypoint.
func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {

	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	const configErrSummary = "Unable to configure client"
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

	const createErrSummary = "Unable to create op client"
	// Create OP cli.
	cli, err := onepasswordcli.NewOpCli(address, email, secretKey)
	if err != nil {
		resp.Diagnostics.AddError(createErrSummary, "Unable to create 1password op cmd client:\n\n"+err.Error())
		return
	}

	// Create  repository.
	repo, err := onepasswordcli.NewRepository(*cli)
	if err != nil {
		resp.Diagnostics.AddError(createErrSummary, "Unable to create 1password op repository:\n\n"+err.Error())
		return
	}

	p.repo = repo
	p.configured = true
}

func (p *provider) configureAddress(config providerData) (addres string, err error) {
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

func (p *provider) configureEmail(config providerData) (addres string, err error) {
	if config.Email.Unknown {
		return "", fmt.Errorf("cannot use unknown value as email")
	}

	// If not set get from env, the value has priority.
	var email string
	if config.Address.Null {
		email = os.Getenv(envVarOpEmail)
	} else {
		email = config.Email.Value
	}

	if email == "" {
		return "", fmt.Errorf("email cannot be an empty string")
	}

	return email, nil
}

func (p *provider) configureSecretKey(config providerData) (addres string, err error) {
	if config.SecretKey.Unknown {
		return "", fmt.Errorf("cannot use unknown value as secret key")
	}

	// If not set get from env, the value has priority.
	var secretKey string
	if config.Address.Null {
		secretKey = os.Getenv(envVarOpSecretKey)
	} else {
		secretKey = config.SecretKey.Value
	}

	if secretKey == "" {
		return "", fmt.Errorf("secret key cannot be an empty string")
	}

	return secretKey, nil
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"onepasswordorg_user":         resourceUserType{},
		"onepasswordorg_group":        resourceGroupType{},
		"onepasswordorg_group_member": resourceGroupMemberType{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}
