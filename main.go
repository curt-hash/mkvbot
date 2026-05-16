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

//go:embed mkvbot.toml
var configBytes []byte

func main() {
	cmd := newCLICommand()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	if cmd.Bool(demoFlagName) {
		return runDemo(ctx)
	}

	if cmd.Bool(createConfigFlagName) {
		if err := os.WriteFile(defaultConfigPath, configBytes, 0600); err != nil {
			return fmt.Errorf("create %q: %w", defaultConfigPath, err)
		}
		fmt.Printf("wrote mkvbot config file %q\n", defaultConfigPath)
		return nil
	}

	if cmd.Bool(createProfileFlagName) {
		if err := os.WriteFile(defaultProfilePath, profileBytes, 0600); err != nil {
			return fmt.Errorf("create %q: %w", defaultProfilePath, err)
		}
		fmt.Printf("wrote makemkv profile %q\n", defaultProfilePath)
		return nil
	}

	profilePath := cmd.String(profileFlagName)
	if _, err := os.Stat(profilePath); err == nil {
		if profilePath, err = filepath.Abs(profilePath); err != nil {
			return fmt.Errorf("get absolute path of %q: %w", profilePath, err)
		}
	} else {
		slog.Warn("profile does not exist", "path", profilePath)
		profilePath = ""
	}

	weights := make(map[string]int64, len(bestTitleHeuristics))
	for _, h := range bestTitleHeuristics {
		weights[h.name] = cmd.Int64(h.flagName)
	}

	cfg := &applicationConfig{
		outputDirPath: cmd.String(outputDirFlagName),
		makemkvConfig: &makemkv.Config{
			ExePath:          cmd.String(makemkvconFlagName),
			ProfilePath:      profilePath,
			ReadCacheSizeMB:  cmd.Int64(cacheFlagName),
			MinLengthSeconds: cmd.Int64(minLengthFlagName),
		},
		debug:                      cmd.Bool(debugFlagName),
		quiet:                      cmd.Bool(quietFlagName),
		bestTitleHeuristicsWeights: weights,
		askForTitle:                cmd.Bool(askForTitleFlagName),
		logFilePath:                cmd.String(logFileFlagName),
	}

	app, err := newApplication(cfg)
	if err != nil {
		return fmt.Errorf("initialize application: %w", err)
	}

	return app.run(ctx)
}
