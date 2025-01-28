package main

import (
	"credo/cmd"
	"credo/config"
	"credo/logger"
	"credo/modules"
	"credo/suggest"
	"credo/version"
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	logger := logger.Get()
	configProvider := config.FileProvider{}
	config, err := configProvider.Get()
	if err != nil {
		logger.Fatal(err)
	}
	modules.RegisterModulesCli(cmd.RootCmd, config)
	cmd.RootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			if err := version.PrintVersion(cmd.Use); err != nil {
				logger.Fatal("Error printing version:", err)
			}
		},
	})
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
