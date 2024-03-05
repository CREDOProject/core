package modules

import (
	"credo/logger"
	"credo/project"
	"fmt"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/cobra"

	goisgiturl "github.com/CREDOProject/go-isgiturl"
)

const gitModuleName = "git"

const gitModuleExample = `
Clone a git repository:
	credo git https://github.com/kendomaniac/rCASC

Clone a git repository at a specific version tag:
	credo git https://github.com/kendomaniac/docker4seq 2.1.2
`

func init() { Register(gitModuleName, func() Module { return &gitModule{} }) }

type gitModule struct{}

func (m *gitModule) Commit(config *Config, result any) error {
	newEntry := result.(gitSpell)
	if Contains(config.Git, newEntry) {
		return ErrAlreadyPresent
	}
	config.Git = append(config.Git, newEntry)
	return nil
}

func (m *gitModule) bareRun(p gitSpell) (gitSpell, error) {
	// Logic to get the latest version or the specified version.
	version := p.Version
	if len(version) == 0 {
		version = "HEAD"
	}

	// Setup a spell entry.
	spell := gitSpell{
		URL:     p.URL,
		Version: version,
	}

	// Try Cloning
	_, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:               p.URL,
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: 1,
	})

	if err != nil {
		return spell, err
	}

	return spell, nil
}

func (m *gitModule) Run(anySpell any) error {
	spell := anySpell.(gitSpell)

	// Obtain the project path
	projectPath, err := project.ProjectPath()
	if err != nil {
		return err
	}
	_, _, _, repoPath := goisgiturl.FindScpLikeComponents(spell.URL)
	joinedPath := path.Join(strings.Split(repoPath, "/")...)
	// Try Clone
	_, err = git.PlainClone(path.Join(*projectPath, gitModuleName, joinedPath), false, &git.CloneOptions{
		URL:               spell.URL,
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *gitModule) BulkRun(config *Config) error {
	for _, gs := range config.Git {
		err := m.Run(gs)
		if err != nil {
			return err
		}
	}
	return nil
}

// Struct containing a Spell Entry for a Git repo.
type gitSpell struct {
	URL     string `yaml:"url"`
	Version string `yaml:"version"`
}

// Function to check equality of two GitSpells
func (s gitSpell) equals(t equatable) bool {
	if o, ok := t.(gitSpell); ok {
		return strings.Compare(s.URL, o.URL) == 0 &&
			strings.Compare(s.Version, o.Version) == 0
	}
	return false
}

func (m *gitModule) CliConfig(conifig *Config) *cobra.Command {
	return &cobra.Command{
		Use:     gitModuleName,
		Short:   "Retrieves a remote git repository.",
		Example: gitModuleExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("%s module requires at least one argument.",
					gitModuleName)
			}
			url := args[0]
			if !goisgiturl.IsGitUrl(url) {
				return fmt.Errorf("\"%s\" doesn't look like a git uri.", url)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			version := ""
			if len(args) > 1 {
				version = args[1]
			}
			spell, err := m.bareRun(gitSpell{
				URL:     args[0],
				Version: version,
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
