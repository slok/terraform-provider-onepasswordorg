package onepasswordcli

import "embed"

var (
	//go:embed op-cli-tfe/*
	//
	// The 1password CLI for Terraform cloud.
	//
	// By default terraform cloud doesn't have op CLI. We embed on the binary built
	// by the architecture used on terraform cloud so we can use it when the app runs.
	//
	// Note: We are embedding linux amd64 op binary in all architecture builds.
	// This will increase in ~10MB the size of all architecture binaries although the
	// binary only can be used in linux amd64. We assume that tradeoffin favor of compiling
	// simplicity.
	EmbeddedOpCli embed.FS
)
