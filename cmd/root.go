package cmd

import (
	"credo/logger"
	"credo/version"

	"github.com/spf13/cobra"
)

// RootCmd is the main command of the Credo application.
// It is a Cobra command that serves as a stub to attach other modules to.
var RootCmd = setup()

func setup() *cobra.Command {
	var base = &cobra.Command{
		Use:   "credo",
		Short: "Credo is a tool for creating reproducible bioinformatics environments.",
		Long: `Credo is a tool for creating reproducible bioinformatics environments.
It helps you manage your dependencies, track changes, and share your work with others.`,
	}
	base.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(_ *cobra.Command, args []string) {
			if err := version.PrintVersion(base.Use); err != nil {
				logger.Get().Fatal("Error printing version:", err)
			}
		},
	})
	return base
}
