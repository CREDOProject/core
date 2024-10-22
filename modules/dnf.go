package modules

import (
	"credo/cache"
	"credo/logger"
	"fmt"
	"strings"

	goosinfo "github.com/CREDOProject/go-osinfo"
	"github.com/spf13/cobra"
)

const dnfModuleName = "dnf"

const dnfModuleShort = "Retrieves a dnf package and its dependencies"

const dnfModuleExample = `
Install a dnf package
	credo dnf vim
`

// Registers the dnfModule.
func init() {
	osinfo, err := goosinfo.Retrieve()
	if err != nil {
		logger.Get().Fatal(err)
	}
	supportedDistributions := map[string]struct{}{
		"rhel":   {},
		"centos": {},
		"fedora": {},
	}
	for _, dist := range osinfo.Like {
		if _, ok := supportedDistributions[dist]; ok {
			Register(dnfModuleName, func() Module { return &dnfModule{} })
			return
		}
	}
	if _, ok := supportedDistributions[osinfo.Distribution]; ok {
		Register(dnfModuleName, func() Module { return &dnfModule{} })
		return
	}
}

// dnfModule is used to manage the dnf scope in the credospell configuration.
type dnfModule struct{}

type dnfSpell struct {
	Name                 string     `yaml:"name"`
	Optional             bool       `yaml:"optional,omitempty"`
	Dependencies         []dnfSpell `yaml:"dependencies,omitempty"`
	ExternalDependencies Config     `yaml:"external_dependencies,omitempty"`
}

// Function used to check if two dnfSpell objects are equal.
func (d dnfSpell) equals(t equatable) bool {
	o, ok := t.(dnfSpell)
	if !ok {
		return false
	}
	if strings.Compare(o.Name, d.Name) != 0 {
		return false
	}
	if o.Optional != d.Optional {
		return false
	}
	if len(o.Dependencies) != len(d.Dependencies) {
		return false
	}
	for i := range o.Dependencies {
		if !d.Dependencies[i].equals(o.Dependencies[i]) {
			return false
		}
	}
	return true
}

// Apply implements Module.
func (d *dnfModule) Apply(any) error {
	panic("unimplemented")
}

// BulkApply implements Module.
func (d *dnfModule) BulkApply(config *Config) error {
	panic("unimplemented")
}

// BulkSave implements Module.
func (d *dnfModule) BulkSave(config *Config) error {
	panic("unimplemented")
}

// CliConfig implements Module.
func (d *dnfModule) CliConfig(config *Config) *cobra.Command {
	return &cobra.Command{
		Args:    d.cobraArgs(),
		Example: dnfModuleExample,
		Run:     d.cobraRun(config),
		Short:   dnfModuleShort,
		Use:     dnfModuleName,
	}
}

// Function used to run the module from the command line.
// It serves as an entry point to the bare run of the dnfModule.
//
// Intended to be used by cobra.
func (d *dnfModule) cobraRun(config *Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		// TODO: implement run
		err := d.Commit(config, nil)
		if err != nil && err != ErrAlreadyPresent {
			logger.Get().Fatal(err)
		}
	}
}

// Function used to validate the arguments passed to the dnf command.
// If no arguments are passed, it returns an error.
// Otherwise it returns nil.
//
// Intended to be used by cobra.
func (*dnfModule) cobraArgs() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s module requires at least one argument.",
				aptModuleName)
		}
		return nil
	}
}

// Commit implements Module.
func (d *dnfModule) Commit(config *Config, result any) error {
	panic("unimplemented")
}

// Save implements Module.
func (d *dnfModule) Save(any) error {
	panic("unimplemented")
}

func (*dnfModule) bareRun(d *dnfSpell) (*dnfSpell, error) {
	if spell := cache.Retrieve(dnfModuleName, d.Name); spell != nil {
		if newSpell, ok := spell.(dnfSpell); ok {
			return &newSpell, nil
		}
	}
	// TODO: Implement bare run
	_ = cache.Insert(dnfModuleName, d.Name, &d)
	return d, nil
}

func (*dnfModule) getDnf() *dnf.Dnf {
	return dnf.New("")
}
