package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	Debug                     bool           `toml:"debug"`
	MakemkvconPath            string         `toml:"makemkvcon_path"`
	ProfilePath               string         `toml:"makemkv_profile_path"`
	CacheSize                 int            `toml:"cache_size"`
	MinLength                 int            `toml:"min_length"`
	OutputDirPath             string         `toml:"output_dir_path"`
	Quiet                     bool           `toml:"quiet"`
	AskForTitle               bool           `toml:"ask_for_title"`
	LogFilePath               string         `toml:"log_file_path"`
	BestTitleHeuristicWeights map[string]int `toml:"best_title_heuristic_weights"`
}

func newDefaultConfig() *config {
	cfg := &config{
		Debug:                     false,
		MakemkvconPath:            "",
		ProfilePath:               defaultProfilePath,
		CacheSize:                 1024,
		MinLength:                 1800,
		OutputDirPath:             ".",
		Quiet:                     false,
		AskForTitle:               false,
		LogFilePath:               "",
		BestTitleHeuristicWeights: make(map[string]int),
	}

	for _, h := range bestTitleHeuristics {
		cfg.BestTitleHeuristicWeights[h.name] = h.weight
	}

	return cfg
}

func newConfigFromFile(path string) (*config, error) {
	c := &config{}
	if _, err := toml.DecodeFile(path, c); err != nil {
		return nil, fmt.Errorf("decode config from %q: %w", path, err)
	}

	return c, nil
}

func (cfg *config) writeToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}
	defer f.Close()

	e := toml.NewEncoder(f)
	if err := e.Encode(cfg); err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	return nil
}
