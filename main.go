package main

import (
	"credo/cmd"
	"credo/config"
	"credo/logger"
	"credo/modules"
	"credo/suggest"
	"fmt"
)

func main() {
	logger := logger.Get()
	configProvider := config.FileProvider{}
	config, err := configProvider.Get()
	if err != nil {
		logger.Fatal(err)
	}
	for _, module := range modules.Modules {
		if moduleConfig := module().CliConfig(config); moduleConfig != nil {
			cmd.RootCmd.AddCommand(moduleConfig)
		}
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
	if err := configProvider.Write(config); err != nil {
		logger.Fatal(err)
	}
	// Print suggestions.
	if suggest.HasSuggestion() {
		fmt.Printf("Package suggestions:\n")
		fmt.Print(suggest.Get().String())
	}
}
