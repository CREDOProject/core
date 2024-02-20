package modules

// Application configuration.
type Config struct {
	Git []GitSpell `yaml:"git,omitempty"`
	Pip []PipSpell `yaml:"pip,omitempty"`
}
