package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/curt-hash/mkvbot/pkg/makemkv"
	"github.com/urfave/cli/v3"
)

const betaKeyURL = "https://forum.makemkv.com/forum/viewtopic.php?t=1053"

var betaKeyRe = regexp.MustCompile(`T-[A-Za-z0-9]{20,}`)

func newUpdateKeyCommand() *cli.Command {
	return &cli.Command{
		Name:  "update-key",
		Usage: "fetch and register the current MakeMKV beta key",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runUpdateKey(ctx, cmd)
		},
	}
}

func runUpdateKey(ctx context.Context, _ *cli.Command) error {
	key, err := fetchBetaKey(ctx)
	if err != nil {
		return fmt.Errorf("fetch beta key: %w", err)
	}

	if err := makemkv.Reg(key); err != nil {
		return fmt.Errorf("register key: %w", err)
	}

	fmt.Printf("registered key %s\n", key)
	return nil
}

func fetchBetaKey(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, betaKeyURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("GET %s: %w", betaKeyURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	key := betaKeyRe.FindString(string(body))
	if key == "" {
		return "", fmt.Errorf("no beta key found on page")
	}

	return key, nil
}
