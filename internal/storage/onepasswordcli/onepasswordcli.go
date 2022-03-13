package onepasswordcli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// OpCli knows how to execute Op CLI commands.
type OpCli interface {
	RunOpCmd(ctx context.Context, args []string) (stdout, stderr string, err error)
}

//go:generate mockery --case underscore --output onepasswordclimock --outpkg onepasswordclimock --name OpCli

type opCli struct {
	binPath      string
	sessionToken string
}

const defaultBinPath = "op"

// NewOpCLI creates a new signed OpCLI command executor.
func NewOpCli(address, email, secretKey, password string) (OpCli, error) {
	binPath := defaultBinPath

	// Login.
	cmd := exec.Command(defaultBinPath, "signin", address, email, secretKey, "--output=raw", "--shorthand=terraform")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		defer stdin.Close()
		_, err := io.WriteString(stdin, fmt.Sprintf("%s\n", password))
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}
	}()

	result, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("cannot signin: %w: %s", err, string(result))
	}

	return opCli{
		binPath:      binPath,
		sessionToken: strings.TrimSpace(string(result)),
	}, nil
}

func (o opCli) RunOpCmd(ctx context.Context, args []string) (stdout, stderr string, err error) {
	if o.sessionToken == "" {
		return "", "", fmt.Errorf("unauthenticated, op cli must singin first")
	}

	// Set session token and account before executing the command.
	args = append([]string{"--session", o.sessionToken, "--account", "terraform"}, args...)

	// Prepare command and execute.
	cmd := exec.CommandContext(ctx, defaultBinPath, args...)
	var sout, serr bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = &serr
	err = cmd.Run()

	return sout.String(), serr.String(), err
}

// NewRepository returns a 1password CLI (op) based respoitory.
func NewRepository(cli OpCli) (*Repository, error) {
	return &Repository{
		cli: cli,
	}, nil
}

// Repository knows how to execute 1password operations using 1password CLI.
type Repository struct {
	cli OpCli
}
