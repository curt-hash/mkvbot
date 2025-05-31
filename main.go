package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/curt-hash/mkvbot/pkg/makemkv"
	"github.com/urfave/cli/v3"
)

//go:embed makemkv.xml
var profileBytes []byte

func main() {
	cmd := newCLICommand()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	if cmd.Bool(createConfigFlagName) {
		c := newDefaultConfig()
		name := defaultConfigPath
		if err := c.writeToFile(defaultConfigPath); err != nil {
			return fmt.Errorf("create %q: %w", name, err)
		}
		fmt.Printf("wrote mkvbot config file %q\n", name)
		return nil
	}

	if cmd.Bool(createProfileFlagName) {
		name := defaultProfilePath
		if err := os.WriteFile(name, profileBytes, 0600); err != nil {
			return fmt.Errorf("create %q: %w", name, err)
		}
		fmt.Printf("wrote makemkv profile %q\n", name)
		return nil
	}

	cfg, err := newConfigFromFile(cmd.String(configFlagName))
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	profilePath := cfg.ProfilePath
	if _, err := os.Stat(profilePath); err == nil {
		if profilePath, err = filepath.Abs(profilePath); err != nil {
			return fmt.Errorf("get absolute path of %q: %w", profilePath, err)
		}
	} else {
		slog.Warn("profile does not exist", "path", profilePath)
		profilePath = ""
	}

	appCfg := &applicationConfig{
		outputDirPath: cfg.OutputDirPath,
		makemkvConfig: &makemkv.Config{
			ExePath:          cfg.MakemkvconPath,
			ProfilePath:      profilePath,
			ReadCacheSizeMB:  int64(cfg.CacheSize),
			MinLengthSeconds: int64(cfg.MinLength),
		},
		debug:                      cfg.Debug,
		quiet:                      cfg.Quiet,
		bestTitleHeuristicsWeights: cfg.BestTitleHeuristicWeights,
		askForTitle:                cfg.AskForTitle,
		logFilePath:                cfg.LogFilePath,
	}

	app, err := newApplication(appCfg)
	if err != nil {
		return fmt.Errorf("initialize application: %w", err)
	}

	return app.run(ctx)
}
