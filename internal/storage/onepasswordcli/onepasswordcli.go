package onepasswordcli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
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

// NewOpCLI creates a new signed OpCLI command executor.
func NewOpCli(customCliPath, address, email, secretKey, password string) (OpCli, error) {
	binPath, err := prepareOpCliBinary(customCliPath)
	if err != nil {
		return nil, fmt.Errorf("could not prepare op cli: %w", err)
	}

	// Login.
	cmd := exec.Command(binPath, "signin", address, email, secretKey, "--output=raw", "--shorthand=terraform")
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

// prepareOpCliBinary will prepare the op binary returning the path the execution must use.
//
// If running outside terraform cloud (tfe), we will require the op tool is available on
// the system path as `op`.
//
// If we are running in terraform cloud, then we will get our op binary from an embedded
// file system, copy to a path and execute from there. We know that tfe uses linux amd64
// machines.
func prepareOpCliBinary(customBinPath string) (binPath string, err error) {
	const (
		defaultBinPath   = "op"
		tfeBinPath       = "/tmp/op-tfe"
		tfeRunningEnvVar = "TFC_RUN_ID"
	)

	// If not terraform cloud, then regular execution.
	tfe := os.Getenv(tfeRunningEnvVar)
	if tfe == "" {
		// If custom binary path is empty, use the path default `op` one.
		if customBinPath == "" {
			return defaultBinPath, nil
		}
		return customBinPath, nil
	}

	// Copy embedded binary into a tmp file.
	f, err := EmbeddedOpCli.ReadFile("op-cli-tfe/op")
	if err != nil {
		return "", fmt.Errorf("could not read embedded op cli: %w", err)
	}

	err = os.WriteFile(tfeBinPath, f, 0755)
	if err != nil {
		return "", fmt.Errorf("could not write embedded op cli into fs: %w", err)
	}

	return tfeBinPath, nil
}

func (o opCli) RunOpCmd(ctx context.Context, args []string) (stdout, stderr string, err error) {
	if o.sessionToken == "" {
		return "", "", fmt.Errorf("unauthenticated, op cli must singin first")
	}

	// Set session token and account before executing the command.
	args = append([]string{"--session", o.sessionToken, "--account", "terraform"}, args...)

	// Prepare command and execute.
	cmd := exec.CommandContext(ctx, o.binPath, args...)
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
