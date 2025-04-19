package main

import (
	"github.com/BurntSushi/toml"
)

type config struct {
	Debug                     bool
	MakemkvconPath            string
	MakemkvProfilePath        string
	CacheSize                 int
	MinLength                 int
	OutputDirPath             string
	Quiet                     bool
	AskForTitle               bool
	LogFilePath               string
	BestTitleHeuristicWeights map[string]int
}

func newDefaultConfig() *config {
	return &config{

	}
}
