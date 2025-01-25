package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/curt-hash/mkvbot/pkg/makemkvcon"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := newCLICommand()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
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
