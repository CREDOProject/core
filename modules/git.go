package modules

import (
	"credo/logger"
	"credo/project"
	"credo/utils"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/cobra"
)

const gitModuleName = "git"

const gitModuleExample = `
Clone a git repository:
	credo git https://github.com/kendomaniac/rCASC

Clone a git repository at a specific version tag:
	credo git https://github.com/kendomaniac/docker4seq 2.1.2
`

type GitModule struct {
	logger *log.Logger
}

func (m *GitModule) Commit(config *Config, result any) error {
	newEntry := result.(GitSpell)

	for _, spell := range config.Git {
		if spell.equals(&newEntry) {
			return nil
		}
	}

	config.Git = append(config.Git, newEntry)
	return nil
}

func (m *GitModule) bareRun(p GitSpell) (GitSpell, error) {
	// Logic to get the latest version or the specified version.
	version := p.Version
	if len(version) == 0 {
		version = "HEAD"
	}

	// Setup a spell entry.
	spell := GitSpell{
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

func (m *GitModule) Run(anySpell any) error {
	spell := anySpell.(GitSpell)

	// Obtain the project path
	projectPath, err := project.ProjectPath()
	if err != nil {
		return err
	}
	_, _, _, repoPath := utils.FindScpLikeComponents(spell.URL)
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

func (m *GitModule) BulkRun(config *Config) error {
	for _, gs := range config.Git {
		err := m.Run(gs)
		if err != nil {
			return err
		}
	}
	return nil
}

// Struct containing a Spell Entry for a Git repo.
type GitSpell struct {
	URL     string `yaml:"url"`
	Version string `yaml:"version"`
}

// Function to check equality of two GitSpells
func (s *GitSpell) equals(t *GitSpell) bool {
	return s.URL == t.URL && s.Version == t.Version
}

func (m *GitModule) Marshaler() interface{} {
	return GitSpell{}
}

func (m *GitModule) CliConfig(conifig *Config) *cobra.Command {
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
			if !utils.IsGitUrl(url) {
				return fmt.Errorf("\"%s\" doesn't look like a git uri.", url)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			version := ""
			if len(args) > 1 {
				version = args[1]
			}
			spell, err := m.bareRun(GitSpell{
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

func init() { Register(gitModuleName, func() Module { return &GitModule{} }) }
