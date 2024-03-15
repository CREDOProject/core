package main

import (
	"credo/cmd"
	"credo/config"
	"credo/logger"
	"credo/modules"
	"credo/storage"

	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	logger := logger.Get()
	store := storage.FileStorage{
		Filename: "credospell.yaml",
	}
	prevFile := store.Read()
	fullConfig, err := config.FromFile(prevFile)
	if err != nil {
		logger.Fatal(err)
	}
	for _, module := range modules.Modules {
		if config := module().CliConfig(&fullConfig); config != nil {
			cmd.RootCmd.AddCommand(config)
		}
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	marshal, err := yaml.Marshal(fullConfig)
	if err != nil {
		logger.Fatal("Can't marshal")
	}

	store.Write(marshal)
}
