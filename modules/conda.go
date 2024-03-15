package modules

import (
	"credo/logger"
	"credo/project"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	goconda "github.com/CREDOProject/go-conda"
	condautils "github.com/CREDOProject/go-conda/utils"
	"github.com/spf13/cobra"
)

const condaModuleName = "conda"

const condaModuleShort = "Retrieves a conda package and its dependencies."

const condaModuleExample = `
Install a conda package:
	credo conda numpy

Install a conda package from a channel:
	credo conda scipy --channel=bioconda
`

// Registers the condaModule.
func init() { Register(condaModuleName, func() Module { return &condaModule{} }) }

// condaModule is used to manage the conda scope in the credospell configuration.
type condaModule struct{}

type condaSpell struct {
	Name    string `yaml:"name"`
	Channel string `yaml:"channel,omitempty"`
}

// Function used to check if two condaSpell objects are equal.
// It takes in an equatable interface as a parameter and returns a boolean
// value indicating whether the two objects are equal or not.
// The function first checks if the input parameter t is of type condaSpell.
//
// If it is, it proceeds to compare the Name and Channel of the two
// objects. The function returns true if the two objects are equal.
// Otherwise, it returns false.
//
// This function is useful for comparing two condaSpell objects to determine if
// they represent the same configuration or not.
func (c condaSpell) equals(t equatable) bool {
	if o, ok := t.(condaSpell); ok {
		return strings.Compare(c.Name, o.Name) == 0 &&
			strings.Compare(c.Channel, o.Channel) == 0
	}
	return false
}

// BulkRun implements Module.
func (c *condaModule) BulkRun(config *Config) error {
	for _, cs := range config.Conda {
		if err := c.Run(cs); err != nil {
			return err
		}
	}
	return nil
}

// Function used to run the module from the command line.
// It serves as an entry point to the bare run of the condaModule.
//
// Intended to be used by cobra.
func (c *condaModule) cobraRun(config *Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		channel, _ := cmd.Flags().GetString("channel")
		spell, err := c.bareRun(condaSpell{
			Name:    args[0],
			Channel: channel,
		})
		if err != nil {
			logger.Get().Fatal(err)
		}
		err = c.Commit(config, spell)
		if err != nil {
			logger.Get().Fatal(err)
		}
	}
}

// Function used to validate the arguments passed to the conda command.
// If no arguments are passed, it returns an error.
// Otherwise it returns nil.
//
// Intended to be used by cobra.
func (m *condaModule) cobraArgs() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s module requires at least one argument.",
				condaModuleName)
		}
		return nil
	}
}

// CliConfig implements Module.
func (c *condaModule) CliConfig(config *Config) *cobra.Command {
	command := &cobra.Command{
		Short:   condaModuleShort,
		Example: condaModuleExample,
		Use:     condaModuleName,
		Run:     c.cobraRun(config),
		Args:    c.cobraArgs(),
	}
	command.PersistentFlags().String("channel", "", "Conda channel to use.")
	return command
}

// Commit implements Module.
func (c *condaModule) Commit(config *Config, result any) error {
	newEntry := result.(condaSpell)
	if Contains(config.Conda, newEntry) {
		return ErrAlreadyPresent
	}
	config.Conda = append(config.Conda, newEntry)
	return nil
}

func (c *condaModule) bareRun(p condaSpell) (condaSpell, error) {
	condaBinary, err := condautils.DetectCondaBinary()
	if err != nil {
		return condaSpell{}, nil
	}
	cmd, err := goconda.New(condaBinary, "", "").
		Install(&goconda.PackageInfo{
			PackageName: p.Name,
			Channel:     p.Channel,
		}).
		DryRun().
		Seal()
	err = cmd.Run(&goconda.RunOptions{
		Output: os.Stdout,
	})
	if err != nil {
		return condaSpell{}, err
	}
	return p, nil
}

// Run implements Module.
func (c *condaModule) Run(anySpell any) error {
	project, err := project.ProjectPath()
	if err != nil {
		return err
	}
	condaBinary, err := condautils.DetectCondaBinary()
	if err != nil {
		return err
	}

	spell, ok := anySpell.(condaSpell)
	if !ok {
		return errors.New("Error converting.")
	}

	downloadPath := path.Join(*project, condaModuleName)
	cmd, err := goconda.
		New(condaBinary, downloadPath, downloadPath).
		Download(&goconda.PackageInfo{
			PackageName: spell.Name,
			Channel:     spell.Channel,
		}, downloadPath).Seal()

	err = cmd.Run(&goconda.RunOptions{
		Output: os.Stdout,
	})

	if err != nil {
		return err
	}
	return nil
}
