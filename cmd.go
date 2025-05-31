package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

var Version string

const (
	configFlagName        = "config"
	createConfigFlagName  = "create-config"
	createProfileFlagName = "create-profile"
)

const (
	defaultConfigPath  = "mkvbot.toml"
	defaultProfilePath = "makemkv.xml"
)

func newCLICommand() *cli.Command {
	cmd := &cli.Command{
		Name:      "mkvbot",
		Version:   Version,
		Usage:     "Automation for makemkv",
		Copyright: "(c) 2025 Curt Hash",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    configFlagName,
				Value:   defaultConfigPath,
				Usage:   "`PATH` to mkvbot config file",
				Aliases: []string{"c"},
			},
			&cli.BoolFlag{
				Name:  createConfigFlagName,
				Usage: fmt.Sprintf("create the default mkvbot config file (%q) then exit", defaultConfigPath),
			},
			&cli.BoolFlag{
				Name:  createProfileFlagName,
				Usage: fmt.Sprintf("create the default makemkv profile file (%q) then exit", defaultProfilePath),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return run(ctx, cmd)
		},
	}

	return cmd
}
