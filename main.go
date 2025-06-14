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

//go:embed profile.xml
var profileBytes []byte

func main() {
	cmd := newCLICommand()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	if cmd.Bool(createProfileFlagName) {
		name := "profile.xml"
		if _, err := os.Stat(name); err == nil {
			return fmt.Errorf("create %q: file exists", name)
		}
		if err := os.WriteFile(name, profileBytes, 0600); err != nil {
			return fmt.Errorf("create %q: %w", name, err)
		}
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
