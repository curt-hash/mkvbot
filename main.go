package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/curt-hash/mkvbot/pkg/makemkvcon"
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
	if cmd.Bool("create-profile") {
		if err := os.WriteFile("profile.xml", profileBytes, 0644); err != nil {
			return fmt.Errorf("create profile.xml: %w", err)
		}
	}

	profilePath := cmd.String("profile")
	if _, err := os.Stat(profilePath); err == nil {
		if profilePath, err = filepath.Abs(profilePath); err != nil {
			return fmt.Errorf("get absolute path of %q: %w", profilePath, err)
		}
	} else {
		slog.Warn("profile does not exist", "path", profilePath)
		profilePath = ""
	}

	cfg := &applicationConfig{
		outputDirPath: cmd.String("output-dir"),
		makemkvConfig: &makemkvcon.MakeMKVConConfig{
			ExePath:          cmd.String("makemkvcon"),
			ProfilePath:      profilePath,
			ReadCacheSizeMB:  cmd.Int("cache"),
			MinLengthSeconds: cmd.Int("minlength"),
		},
	}

	app, err := newApplication(cfg)
	if err != nil {
		return fmt.Errorf("initialize application: %w", err)
	}

	return app.run(ctx)
}
