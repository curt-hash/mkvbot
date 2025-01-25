package main

import (
	"context"

	"github.com/urfave/cli/v3"
)

var Version string

func newCLICommand() *cli.Command {
	cmd := &cli.Command{
		Name:      "mkvbot",
		Version:   Version,
		Usage:     "Automation for makemkv",
		Copyright: "(c) 2025 Curt Hash",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "log debug messages",
			},
			&cli.StringFlag{
				Name:    "makemkvcon",
				Value:   "",
				Usage:   "`PATH` to makemkvcon executable",
				Aliases: []string{"m"},
			},
			&cli.StringFlag{
				Name:    "profile",
				Value:   "profile.xml",
				Usage:   "pass --profile=`PATH` to makemkv",
				Aliases: []string{"p"},
			},
			&cli.BoolFlag{
				Name:  "create-profile",
				Usage: "create a default profile.xml for use with --profile",
			},
			&cli.IntFlag{
				Name:    "cache",
				Value:   1024,
				Usage:   "pass --cache=`SIZE` to makemkv",
				Aliases: []string{"c"},
			},
			&cli.IntFlag{
				Name:    "minlength",
				Value:   1800,
				Usage:   "pass --minlength=`N` to makemkv",
				Aliases: []string{"l"},
			},
			&cli.StringFlag{
				Name:    "output-dir",
				Value:   ".",
				Usage:   "create output files in `DIR`",
				Aliases: []string{"o"},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return run(ctx, cmd)
		},
	}

	for _, h := range bestTitleHeuristics {
		cmd.Flags = append(cmd.Flags, &cli.IntFlag{
			Name:        h.flagName,
			Value:       h.weight,
			Usage:       h.flagUsage,
			Destination: &h.weight,
		})
	}

	return cmd
}
