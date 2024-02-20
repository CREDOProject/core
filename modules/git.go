package modules

import (
	"errors"
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GitModule struct {
	logger *log.Logger
}

func (m *GitModule) Commit(config *Config, result any) error {
	newEntry := result.(GitSpell) // Type conversion.

	for _, spell := range config.Git {
		if spell.equals(&newEntry) {
			return nil // Break from the for loop.
		}
	}

	config.Git = append(config.Git, newEntry)
	return nil
}

func (m *GitModule) BareRun(c *Config, p *Parameters) any {
	spell, err := m.bareRun(p)
	if err != nil {
		m.logger.Fatal(err)
	}
	return spell
}

func (m *GitModule) bareRun(p *Parameters) (GitSpell, error) {
	if len(p.Env) < 1 {
		return GitSpell{}, errors.New("Git module requires at least one parameter.")
	}
	// Logic to get the latest version or the specified version.
	version := p.Env["version"]
	if len(version) == 0 {
		version = "HEAD"
	}

	// Setup a spell entry.
	spell := GitSpell{
		URL:     p.Env["url"],
		Version: version,
	}

	// Try Cloning
	_, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:               p.Env["url"],
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: 1,
		Tags:              0,
	})

	if err != nil {
		return spell, err
	}

	return spell, nil
}

func (m *GitModule) Run(anySpell any) error {
	spell := anySpell.(GitSpell)

	// Try Clone
	_, err := git.PlainClone("/tmp/test", false, &git.CloneOptions{
		URL:               spell.URL,
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: 1,
		Tags:              0,
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

func init() {
	Register("git", func() Module {
		return &GitModule{}
	})
}
