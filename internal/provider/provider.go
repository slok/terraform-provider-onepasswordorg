package provider

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

type ProviderConfig struct {
	configured bool
	repo       storage.Repository
}

// Provider configuration.
type providerData struct {
	Address         string
	Email           string
	SecretKey       string
	Password        string
	FakeStoragePath string
	CliPath         string
}

func (p *ProviderConfig) configureAddress(config providerData) (string, error) {

	// If not set get from env, the value has priority.
	var address string
	if config.Address == "" {
		address = os.Getenv(envVarOpAddress)
	} else {
		address = config.Address
	}

	if address == "" {
		return "", fmt.Errorf("username cannot be an empty string")
	}

	return address, nil
}

func (p *ProviderConfig) configureEmail(config providerData) (string, error) {

	// If not set get from env, the value has priority.
	var email string
	if config.Email == "" {
		email = os.Getenv(envVarOpEmail)
	} else {
		email = config.Email
	}

	if email == "" {
		return "", fmt.Errorf("email cannot be an empty string")
	}

	return email, nil
}

func (p *ProviderConfig) configureSecretKey(config providerData) (string, error) {

	// If not set get from env, the value has priority.
	var secretKey string
	if config.SecretKey == "" {
		secretKey = os.Getenv(envVarOpSecretKey)
	} else {
		secretKey = config.SecretKey
	}

	if secretKey == "" {
		return "", fmt.Errorf("secret key cannot be an empty string")
	}

	return secretKey, nil
}

func (p *ProviderConfig) configurePassword(config providerData) (string, error) {

	// If not set get from env, the value has priority.
	var password string
	if config.Password == "" {
		password = os.Getenv(envVarOpPassword)
	} else {
		password = config.Password
	}

	if password == "" {
		return "", fmt.Errorf("password cannot be an empty string")
	}

	return password, nil
}

func (p *ProviderConfig) configureFakeStoragePath(config providerData) (string, error) {
	// If not set get from env, the value has priority.
	var fakePath string
	if config.FakeStoragePath == "" {
		fakePath = os.Getenv(EnvVarOpFakeStoragePath)
	} else {
		fakePath = config.FakeStoragePath
	}

	return fakePath, nil
}

func (p *ProviderConfig) configureCliPath(config providerData) (string, error) {
	// If not set get from env, the value has priority.
	var cliPath string
	if config.FakeStoragePath == "" {
		cliPath = os.Getenv(EnvVarOpCliPath)
	} else {
		cliPath = config.CliPath
	}

	return cliPath, nil
}

// Provider The 1Password Connect terraform provider
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: fmt.Sprintf("Set account 1password domain address (e.g: something.1password.com). Also `%s` env var can be used.", envVarOpAddress),
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: fmt.Sprintf("Set account 1password email. Also `%s` env var can be used.", envVarOpEmail),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: fmt.Sprintf("Set account 1password secret key. Also `%s` env var can be used.", envVarOpSecretKey),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: fmt.Sprintf("Set account 1password password. Also `%s` env var can be used.", envVarOpPassword),
			},
			"fake_storage_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: fmt.Sprintf("File to a path where the provider will store the data as if it is 1password (this is used only on development). Also `%s` env var can be used.", EnvVarOpFakeStoragePath),
			},
			"op_cli_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: fmt.Sprintf("The path that points to the op cli binary. Also `%s` env var can be used. (by default `op` on system path, ignored if run in Terraform cloud).", EnvVarOpCliPath),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"onepasswordorg_group": dataSourceGroup(),
			"onepasswordorg_item":  dataSourceItem(),
			"onepasswordorg_user":  dataSourceUser(),
			"onepasswordorg_vault": dataSourceVault(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"onepasswordorg_group":              resourceGroup(),
			"onepasswordorg_group_member":       resourceGroupMember(),
			"onepasswordorg_item":               resourceItem(),
			"onepasswordorg_user":               resourceUser(),
			"onepasswordorg_vault":              resourceVault(),
			"onepasswordorg_vault_group_access": resourceVaultGroupAccess(),
			"onepasswordorg_vault_user_access":  resourceVaultUserAccess(),
		},
	}
	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		p := ProviderConfig{}
		config := providerData{
			Address:         d.Get("address").(string),
			Email:           d.Get("email").(string),
			SecretKey:       d.Get("secret_key").(string),
			Password:        d.Get("password").(string),
			FakeStoragePath: d.Get("fake_storage_path").(string),
			CliPath:         d.Get("op_cli_path").(string),
		}
		// Error summaries
		const (
			configErrSummary = "Unable to configure client:"
			createErrSummary = "Unable to create op client:"
		)

		// Get if we are in fake mode.
		fakeStoragePath, err := p.configureFakeStoragePath(config)
		if err != nil {
			return nil, fmt.Errorf(configErrSummary + "Invalid fake storage path:\n\n" + err.Error())
		}

		// Create fake or regular mode.
		// If the user has set the fake storage path then we are going to use a fake repository.
		// If the user didn't, we will use the op cli based repository (a.k.a real 1password APIs).
		var repo storage.Repository
		if fakeStoragePath != "" {
			repo, err = fake.NewRepository(fakeStoragePath)
			if err != nil {
				return nil, fmt.Errorf(createErrSummary + "Unable to create 1password fake storage:\n\n" + err.Error())
			}
		} else {
			address, err := p.configureAddress(config)
			if err != nil {
				return nil, fmt.Errorf(configErrSummary + "Invalid address:\n\n" + err.Error())
			}

			email, err := p.configureEmail(config)
			if err != nil {
				return nil, fmt.Errorf(configErrSummary + "Invalid email:\n\n" + err.Error())
			}

			secretKey, err := p.configureSecretKey(config)
			if err != nil {
				return nil, fmt.Errorf(configErrSummary + "Invalid secret key:\n\n" + err.Error())
			}

			password, err := p.configurePassword(config)
			if err != nil {
				return nil, fmt.Errorf(configErrSummary + "Invalid password:\n\n" + err.Error())
			}

			cliPath, err := p.configureCliPath(config)
			if err != nil {
				return nil, fmt.Errorf(configErrSummary + "Invalid cli path:\n\n" + err.Error())
			}

			// Create OP cli.
			cli, err := onepasswordcli.NewOpCli(cliPath, address, email, secretKey, password)
			if err != nil {
				return nil, fmt.Errorf(createErrSummary + "Unable to create 1password op cmd client:\n\n" + err.Error())
			}

			// Create  repository.
			repo, err = onepasswordcli.NewRepository(cli)
			if err != nil {
				return nil, fmt.Errorf(createErrSummary + "Unable to create 1password op repository:\n\n" + err.Error())
			}
		}
		p.repo = repo
		p.configured = true

		return p, nil
	}
	return provider
}
