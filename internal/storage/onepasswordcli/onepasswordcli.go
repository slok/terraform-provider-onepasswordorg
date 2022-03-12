package onepasswordcli

// OpCli knows how to execute Op CLI commands.
type OpCli struct{}

// NewOpCLI creates a new signed OpCLI command executor.
func NewOpCli(address, email, secretKey string) (*OpCli, error) {
	return &OpCli{}, nil
}

// NewRepository returns a 1password CLI (op) based respoitory.
func NewRepository(cli OpCli) (*Repository, error) {
	return &Repository{}, nil
}

// Repository knows how to execute 1password operations using 1password CLI.
type Repository struct{}
