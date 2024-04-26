package modules

import (
	"credo/logger"
	"credo/project"
	"credo/suggest"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/CREDOProject/go-apt-client"
	goosinfo "github.com/CREDOProject/go-osinfo"
	"github.com/spf13/cobra"
)

const aptModuleName = "apt"

const aptModuleShort = "Retrieves an apt package and its depenencies."

const aptModuleExample = `
Install a apt package:
	credo apt python3
`

var isAptOptional = regexp.MustCompile(`\<(?P<name>..*)\>`)

// Registers the aptModule.
func init() {
	osinfo, err := goosinfo.Retrieve()
	if err != nil {
		logger.Get().Fatal(err)
	}
	supportedDistributions := map[string]struct{}{
		"ubuntu": {},
		"debian": {},
	}
	if _, ok := supportedDistributions[osinfo.Distribution]; !ok {
		return
	}
	Register(aptModuleName, func() Module { return &aptModule{} })
}

// aptModule is used to manage the apt scope in the credospell configuration.
type aptModule struct{}

// Apply implements Module.
func (m *aptModule) Apply(any) error {
	panic("unimplemented")
}

// BulkApply implements Module.
func (m *aptModule) BulkApply(config *Config) error {
	panic("unimplemented")
}

type aptSpell struct {
	Name                 string     `yaml:"name"`
	Optional             bool       `yaml:"optional,omitempty"`
	Depencencies         []aptSpell `yaml:"dependencies,omitempty"`
	ExternalDependencies Config     `yaml:"external_dependencies,omitempty"`
}

// Function used to check if two aptSpell objects are equal.
// It takes in an equatable interface as a parameter and returns a boolean
// value indicating whether the two objects are equal or not.
// The function first checks if the input parameter t is of type aptSpell.
//
// If it is, it proceeds to compare the Name and Optional of the two
// objects and all its other Depencencies.
// The function returns true if the two objects are equal.
// Otherwise, it returns false.
func (a aptSpell) equals(t equatable) bool {
	if o, ok := t.(aptSpell); ok {
		equality := len(o.Depencencies) == len(a.Depencencies)
		if !equality {
			return false
		}
		for i := range o.Depencencies {
			equality = equality &&
				o.Depencencies[i].equals(a.Depencencies[i])
		}
		return equality
	}
	return false
}

// BulkSave implements Module.
func (m *aptModule) BulkSave(config *Config) error {
	for _, as := range config.Apt {
		for _, dep := range as.Depencencies {
			err := m.Save(dep)
			if err != nil {
				return err
			}
		}
		err := m.Save(as)
		if err != nil {
			return err
		}
	}
	return nil
}

// CliConfig implements Module.
func (m *aptModule) CliConfig(config *Config) *cobra.Command {
	return &cobra.Command{
		Args:    m.cobraArgs(),
		Example: aptModuleExample,
		Run:     m.cobraRun(config),
		Short:   aptModuleShort,
		Use:     aptModuleName,
	}
}

// Function used to validate the arguments passed to the apt command.
// If no arguments are passed, it returns an error.
// Otherwise it returns nil.
//
// Intended to be used by cobra.
func (m *aptModule) cobraArgs() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s module requires at least one argument.",
				aptModuleName)
		}
		return nil
	}
}

// Function used to run the module from the command line.
// It serves as an entry point to the bare run of the aptModule.
//
// Intended to be used by cobra.
func (m *aptModule) cobraRun(config *Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		name := args[0]
		spell, err := m.bareRun(aptSpell{
			Name: name,
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

func (*aptModule) bareRun(spell aptSpell) (aptSpell, error) {
	aptPack := &apt.Package{
		Name: spell.Name,
	}
	output, err := apt.InstallDry(aptPack)
	logger.Get().Print(string(output))
	if err != nil {
		return aptSpell{}, err
	}
	depList, err := apt.GetDependencies(aptPack)
	if err != nil {
		return aptSpell{}, err
	}
	for _, dependency := range depList {
		isOptional := isAptOptional.MatchString(dependency)
		cleanDependency := dependency
		matches := isAptOptional.FindStringSubmatch(dependency)
		nameIndex := isAptOptional.SubexpIndex("name")
		if nameIndex != -1 && isOptional {
			cleanDependency = matches[nameIndex]
			suggest.Register(suggest.Suggestion{
				Module:    aptModuleName,
				From:      aptPack.Name,
				Suggested: cleanDependency,
			})
		}
		spell.Depencencies = append(spell.Depencencies, aptSpell{
			Name:     cleanDependency,
			Optional: isOptional,
		})
	}
	return spell, nil
}

// Commit implements Module.
func (*aptModule) Commit(config *Config, result any) error {
	newEntry, ok := result.(aptSpell)
	if !ok {
		return ErrConverting
	}
	if Contains(config.Apt, newEntry) {
		return ErrAlreadyPresent
	}
	config.Apt = append(config.Apt, newEntry)
	return nil
}

// Save implements Module.
func (*aptModule) Save(anySpell any) error {
	spell, ok := anySpell.(aptSpell)
	if !ok {
		return ErrConverting
	}
	project, err := project.ProjectPath()
	if err != nil {
		return err
	}
	downloadPath := path.Join(*project, aptModuleName)
	os.MkdirAll(downloadPath, 0755)
	aptPack := &apt.Package{
		Name: spell.Name,
	}
	out, err := apt.Download(aptPack, downloadPath)
	logger.Get().Print(string(out))
	return err
}
