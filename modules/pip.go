package modules

import (
	"credo/logger"
	"credo/project"
	"fmt"
	"os"
	"path"
	"strings"

	gopip "github.com/CREDOProject/go-pip"
	"github.com/CREDOProject/go-pip/utils"
	pythonvenv "github.com/CREDOProject/go-pythonvenv"
	"github.com/spf13/cobra"
)

const pipModuleName = "pip"

const pipModuleExample = `
Install a pip package:
	credo pip numpy

Install a pip package pinning it to a version:
	credo pip numpy==1.26.0
`

func init() { Register(pipModuleName, func() Module { return &pipModule{} }) }

type pipModule struct{}

type pipSpell struct {
	Name string `yaml:"name"`
}

// Function to check equality of two PipSpells
func (s pipSpell) equals(t equatable) bool {
	if o, ok := t.(pipSpell); ok {
		return strings.Compare(s.Name, o.Name) == 0
	}
	return false
}

func (m *pipModule) Commit(config *Config, result any) error {
	newEntry := result.(pipSpell)
	if Contains(config.Pip, newEntry) {
		return ErrAlreadyPresent
	}
	config.Pip = append(config.Pip, newEntry)
	return nil
}

func setupPythonVenv(path string) (string, error) {
	venv, err := pythonvenv.Create(path)
	if err != nil {
		return "", err
	}
	return venv.Path, nil
}

func getPipBinary() (*string, error) {
	// Obtain the project path
	projectPath, err := project.ProjectPath()
	if err != nil {
		return nil, err
	}

	venvPath, err := setupPythonVenv(path.Join(*projectPath, "venv"))
	if err != nil {
		return nil, err
	}

	pipBinary, err := utils.PipBinaryFrom(path.Join(venvPath, "bin"))
	if err != nil {
		return nil, err
	}
	return &pipBinary, nil
}

func (m *pipModule) bareRun(p pipSpell) (pipSpell, error) {
	pipBinary, err := getPipBinary()
	if err != nil {
		return pipSpell{}, err
	}

	cmd, err := gopip.New(*pipBinary).Install(p.Name).DryRun().Seal()
	if err != nil {
		return pipSpell{}, err
	}

	err = cmd.Run(&gopip.RunOptions{
		Output: os.Stdout,
	})
	if err != nil {
		return pipSpell{}, err
	}

	return p, nil
}

func (m *pipModule) Run(anySpell any) error {
	project, err := project.ProjectPath()
	if err != nil {
		return err
	}
	pipBinary, err := getPipBinary()
	if err != nil {
		return err
	}
	downloadPath := path.Join(*project, pipModuleName)
	cmd, err := gopip.New(*pipBinary).
		Download(anySpell.(pipSpell).Name, downloadPath).
		Seal()
	err = cmd.Run(&gopip.RunOptions{
		Output: os.Stdout,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *pipModule) BulkRun(config *Config) error {
	for _, ps := range config.Pip {
		err := m.Run(ps)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *pipModule) CliConfig(conifig *Config) *cobra.Command {
	return &cobra.Command{
		Use:     pipModuleName,
		Short:   "Retrieves a python package.",
		Example: pipModuleExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("%s module requires at least one argument.",
					pipModuleName)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			spell, err := m.bareRun(pipSpell{
				Name: args[0],
			})
			if err != nil {
				logger.Get().Fatal(err)
			}
			err = m.Commit(conifig, spell)
			if err != nil {
				logger.Get().Fatal(err)
			}
		},
	}
}
