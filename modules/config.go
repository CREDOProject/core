package modules

// Application configuration.
type Config struct {
	Git []GitSpell `yaml:"git"`
	Pip []PipSpell `yaml:"pip"`
}
