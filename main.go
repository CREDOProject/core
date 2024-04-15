package main

import (
	"credo/cmd"
	"credo/config"
	"credo/modules"

	"fmt"
	"os"
)

func main() {
	configProvider := config.FileProvider{}
	config, err := configProvider.Get()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, module := range modules.Modules {
		if moduleConfig := module().CliConfig(config); moduleConfig != nil {
			cmd.RootCmd.AddCommand(moduleConfig)
		}
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := configProvider.Write(config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
