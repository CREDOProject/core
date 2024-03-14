package modules

import (
	"credo/logger"
	"os"
	"strings"

	goconda "github.com/CREDOProject/go-conda"
	condautils "github.com/CREDOProject/go-conda/utils"
	"github.com/spf13/cobra"
)

const condaModuleName = "conda"

func init() { Register(condaModuleName, func() Module { return &condaModule{} }) }

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

// CliConfig implements Module.
func (c *condaModule) CliConfig(config *Config) *cobra.Command {
	return &cobra.Command{
		Use: condaModuleName,
		Run: func(cmd *cobra.Command, args []string) {
			spell, err := c.bareRun(condaSpell{
				Name: args[0],
			})
			if err != nil {
				logger.Get().Fatal(err)
			}
			err = c.Commit(config, spell)
			if err != nil {
				logger.Get().Fatal(err)
			}
		},
	}
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
func (c *condaModule) Run(any) error {
	return nil
}
