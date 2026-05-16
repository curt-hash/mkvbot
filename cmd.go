package main

import (
	"context"
	"fmt"

	altsrc "github.com/urfave/cli-altsrc/v3"
	"github.com/urfave/cli-altsrc/v3/toml"
	"github.com/urfave/cli/v3"
)

var Version string

const (
	configFlagName        = "config"
	createConfigFlagName  = "create-config"
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

const (
	defaultConfigPath  = "mkvbot.toml"
	defaultProfilePath = "makemkv.xml"
)

// configPath holds the resolved --config value. It is populated when urfave/cli
// parses the --config flag (via Destination) and is read indirectly by every
// other flag's TOML source through a StringPtrSourcer, which dereferences it
// lazily at Lookup time.
var configPath string

func newCLICommand() *cli.Command {
	configSource := altsrc.NewStringPtrSourcer(&configPath)

	tomlSource := func(key string) cli.ValueSourceChain {
		return cli.NewValueSourceChain(toml.TOML(key, configSource))
	}

	cmd := &cli.Command{
		Name:      "mkvbot",
		Version:   Version,
		Usage:     "Automation for makemkv",
		Copyright: "(c) 2025–2026 Curt Hash",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        configFlagName,
				Value:       defaultConfigPath,
				Usage:       "`PATH` to mkvbot config file (TOML)",
				Destination: &configPath,
			},
			&cli.BoolFlag{
				Name:  createConfigFlagName,
				Usage: fmt.Sprintf("write a default mkvbot config file (%q) then exit", defaultConfigPath),
			},
			&cli.BoolFlag{
				Name:    debugFlagName,
				Value:   false,
				Usage:   "log debug messages",
				Sources: tomlSource("debug"),
			},
			&cli.StringFlag{
				Name:    makemkvconFlagName,
				Value:   "",
				Usage:   "`PATH` to makemkvcon executable",
				Aliases: []string{"m"},
				Sources: tomlSource("makemkvcon_path"),
			},
			&cli.StringFlag{
				Name:    profileFlagName,
				Value:   defaultProfilePath,
				Usage:   "pass --profile=`PATH` to makemkv",
				Aliases: []string{"p"},
				Sources: tomlSource("makemkv_profile_path"),
			},
			&cli.BoolFlag{
				Name:  createProfileFlagName,
				Usage: fmt.Sprintf("write the default makemkv profile (%q) then exit", defaultProfilePath),
			},
			&cli.Int64Flag{
				Name:    cacheFlagName,
				Value:   1024,
				Usage:   "pass --cache=`SIZE`MiB to makemkv",
				Aliases: []string{"c"},
				Sources: tomlSource("cache_size"),
			},
			&cli.Int64Flag{
				Name:    minLengthFlagName,
				Value:   1800,
				Usage:   "pass --minlength=`N` to makemkv",
				Aliases: []string{"l"},
				Sources: tomlSource("min_length"),
			},
			&cli.StringFlag{
				Name:    outputDirFlagName,
				Value:   ".",
				Usage:   "create output files in `DIR`",
				Aliases: []string{"o"},
				Sources: tomlSource("output_dir_path"),
			},
			&cli.BoolFlag{
				Name:    quietFlagName,
				Usage:   "do not beep",
				Aliases: []string{"q"},
				Sources: tomlSource("quiet"),
			},
			&cli.BoolFlag{
				Name:    askForTitleFlagName,
				Usage:   "ask you to choose the best title",
				Aliases: []string{"a"},
				Sources: tomlSource("ask_for_title"),
			},
			&cli.StringFlag{
				Name:    logFileFlagName,
				Usage:   "append log messages to `FILE`",
				Aliases: []string{"L"},
				Sources: tomlSource("log_file_path"),
			},
		},
		Commands: []*cli.Command{
			newUpdateKeyCommand(),
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
			Sources:     tomlSource("best_title_heuristic_weights." + h.name),
		})
	}

	return cmd
}
