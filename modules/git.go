package modules

import (
	"credo/logger"
	"credo/project"
	"fmt"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/cobra"

	goisgiturl "github.com/CREDOProject/go-isgiturl"
)

const gitModuleName = "git"

const gitModuleShort = "Retrieves a remote git repository."

const gitModuleExample = `
Clone a git repository:
	credo git https://github.com/kendomaniac/rCASC

Clone a git repository at a specific version tag:
	credo git https://github.com/kendomaniac/docker4seq 2.1.2
`

// Registers the gitModule.
func init() { Register(gitModuleName, func() Module { return &gitModule{} }) }

// gitModule is used to manage the git scope in the credospell configuration.
type gitModule struct{}

func (m *gitModule) Commit(config *Config, result any) error {
	newEntry, ok := result.(gitSpell)
	if !ok {
		return ErrConverting
	}
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
		ReferenceName:     plumbing.NewBranchReferenceName(version),
	})

	if err != nil {
		return spell, err
	}

	return spell, nil
}

func (m *gitModule) Run(anySpell any) error {
	spell, ok := anySpell.(gitSpell)
	if !ok {
		return ErrConverting
	}

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
		ReferenceName:     plumbing.NewBranchReferenceName(spell.Version),
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
	URL                  string `yaml:"url"`
	Version              string `yaml:"version"`
	ExternalDependencies Config `yaml:"external_dependencies,omitempty"`
}

// Function used to check if two aptSpell objects are equal.
// It takes in an equatable interface as a parameter and returns a boolean
// value indicating whether the two objects are equal or not.
// The function first checks if the input parameter t is of type gitSpell.
//
// If it is, it proceeds to compare the URL and Version of the two
// objects.
// The function returns true if the two objects are equal.
// Otherwise, it returns false.
func (s gitSpell) equals(t equatable) bool {
	if o, ok := t.(gitSpell); ok {
		return strings.Compare(s.URL, o.URL) == 0 &&
			strings.Compare(s.Version, o.Version) == 0
	}
	return false
}

// Function used to run the module from the command line.
// It serves as an entry point to the bare run of the gitModule.
//
// Intended to be used by cobra.
func (m *gitModule) cobraRun(config *Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
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
		err = m.Commit(config, spell)
		if err != nil {
			logger.Get().Fatal(err)
		}
	}
}

// Function used to validate the arguments passed to the git command.
// If no arguments are passed, it returns an error.
// Otherwise it returns nil.
//
// Intended to be used by cobra.
func (m *gitModule) cobraArgs() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s module requires at least one argument.",
				gitModuleName)
		}
		url := args[0]
		if !goisgiturl.IsGitUrl(url) {
			return fmt.Errorf("\"%s\" doesn't look like a git uri.", url)
		}
		return nil
	}
}

// CliConfig implements Module.
func (m *gitModule) CliConfig(config *Config) *cobra.Command {
	return &cobra.Command{
		Use:     gitModuleName,
		Short:   gitModuleShort,
		Example: gitModuleExample,
		Args:    m.cobraArgs(),
		Run:     m.cobraRun(config),
	}
}
