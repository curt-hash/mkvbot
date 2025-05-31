package main

import (
	"context"

	"github.com/urfave/cli/v3"
)

var Version string

const (
	debugFlagName         = "debug"
	makemkvconFlagName    = "makemkvcon"
	profileFlagName       = "profile"
	createProfileFlagName = "create-profile"
	cacheFlagName         = "cache"
	minLengthFlagName     = "minlength"
	outputDirFlagName     = "output-dir"
	quietFlagName         = "quiet"
	askForTitleFlagName   = "ask-title"
	logFileFlagName       = "log"
)

func newCLICommand() *cli.Command {
	cmd := &cli.Command{
		Name:      "mkvbot",
		Version:   Version,
		Usage:     "Automation for makemkv",
		Copyright: "(c) 2025 Curt Hash",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  debugFlagName,
				Value: false,
				Usage: "log debug messages",
			},
			&cli.StringFlag{
				Name:    makemkvconFlagName,
				Value:   "",
				Usage:   "`PATH` to makemkvcon executable",
				Aliases: []string{"m"},
			},
			&cli.StringFlag{
				Name:    profileFlagName,
				Value:   "profile.xml",
				Usage:   "pass --profile=`PATH` to makemkv",
				Aliases: []string{"p"},
			},
			&cli.BoolFlag{
				Name:  createProfileFlagName,
				Usage: "create a default profile.xml for use with --profile",
			},
			&cli.IntFlag{
				Name:    cacheFlagName,
				Value:   1024,
				Usage:   "pass --cache=`SIZE` to makemkv",
				Aliases: []string{"c"},
			},
			&cli.IntFlag{
				Name:    minLengthFlagName,
				Value:   1800,
				Usage:   "pass --minlength=`N` to makemkv",
				Aliases: []string{"l"},
			},
			&cli.StringFlag{
				Name:    outputDirFlagName,
				Value:   ".",
				Usage:   "create output files in `DIR`",
				Aliases: []string{"o"},
			},
			&cli.BoolFlag{
				Name:    quietFlagName,
				Usage:   "do not beep",
				Aliases: []string{"q"},
			},
			&cli.BoolFlag{
				Name:    askForTitleFlagName,
				Usage:   "ask you to choose the best title",
				Aliases: []string{"a"},
			},
			&cli.StringFlag{
				Name:    logFileFlagName,
				Usage:   "append log messages to `FILE`",
				Aliases: []string{"L"},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return run(ctx, cmd)
		},
	}

	for _, h := range bestTitleHeuristics {
		cmd.Flags = append(cmd.Flags, &cli.Int64Flag{
			Name:        h.flagName,
			Value:       h.weight,
			Usage:       h.flagUsage,
			Destination: &h.weight,
		})
	}

	return cmd
}
