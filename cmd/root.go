package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd is the main command of the Credo application.
// It is a Cobra command that serves as a stub to attach other modules to.
var RootCmd = &cobra.Command{
	Use:   "credo",
	Short: "Credo is a tool for creating reproducible bioinformatics environments.",
	Long: `Credo is a tool for creating reproducible bioinformatics environments.
It helps you manage your dependencies, track changes, and share your work with others.`,
}
