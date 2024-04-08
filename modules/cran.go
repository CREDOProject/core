package modules

import (
	"credo/logger"
	"credo/project"
	"fmt"
	"os"
	"path"
	"strings"

	rcran "github.com/CREDOProject/go-rcran"
	rscript "github.com/CREDOProject/go-rscript"
	"github.com/spf13/cobra"
)

const cranModuleName = "cran"
const bioconductorModuleName = "bioconductor"

const cranModuleShort = "Retrieves a CRAN package and its dependencies."

const cranModuleExample = `
Install a package from CRAN
	credo cran abind

Install a package from BioConductor
	credo bioconductor GenomicRanges
`

// Registers the carnModule.
func init() { Register(cranModuleName, func() Module { return &cranModule{} }) }

// cranModule is used to manage the CARN scope in the credospell configuration.
type cranModule struct{}

type cranSpell struct {
	PackageName  string      `yaml:"package_name,omitempty"`
	PackagePath  string      `yaml:"package_path,omitempty"`
	Repository   string      `yaml:"repository,omitempty"`
	BioConductor bool        `yaml:"bioconductor,omitempty"`
	Dependencies []cranSpell `yaml:"dependencies,omitempty"`
}

// spellFromDownloadOptions returns a new cran spell from *rcran.DownloadOptions.
func (*cranModule) spellFromDownloadOptions(
	options *rcran.DownloadOptions,
) *cranSpell {
	spell := &cranSpell{
		PackageName: options.PackageName,
		Repository:  options.Repository,
	}
	return spell
}

// equals checks if two cranSpell objects are equal.
func (c cranSpell) equals(t equatable) bool {
	// TODO: implement equality check.
	s, ok := t.(cranSpell)
	if !ok {
		return false
	}
	equality := len(s.Dependencies) == len(c.Dependencies)
	if !equality {
		return false
	}
	for i := range s.Dependencies {
		equality = equality &&
			s.Dependencies[i].equals(c.Dependencies[i])
	}
	return equality && strings.Compare(s.PackageName, c.PackageName) == 0 &&
		strings.Compare(s.PackagePath, c.PackagePath) == 0 &&
		s.BioConductor == c.BioConductor
}

// BulkRun implements Module.
func (c *cranModule) BulkRun(config *Config) error {
	for _, cs := range config.Cran {
		if err := c.Run(cs); err != nil {
			return err
		}
	}
	return nil
}

func (m *cranModule) bareRun(c cranSpell, cfg *Config) (*cranSpell, error) {
	if c.BioConductor {
		err := m.installBioconductor(cfg)
		if err != nil {
			return nil, err
		}
	}
	bin, err := rscript.DetectRscriptBinary()
	if err != nil {
		return nil, err
	}
	finalSpell, err := m.bareRunSingle(c, bin, c.BioConductor)
	// Retrieve dependencies
	cmd := ""
	if c.BioConductor {
		cmd, err = rcran.GetBioconductorDepenencies(&rcran.InstallOptions{
			PackageName: c.PackageName,
			Repository:  c.Repository,
			DryRun:      true,
		})
	} else {
		cmd, err = rcran.GetDependencies(&rcran.InstallOptions{
			PackageName: c.PackageName,
			Repository:  c.Repository,
			DryRun:      true,
		})
	}
	if err != nil {
		return nil, err
	}
	script, err := rscript.New(bin).Evaluate(cmd).Seal()
	if err != nil {
		return nil, err
	}
	out, err := script.CombinedOutput()
	outString := string(out)
	if err != nil {
		return nil, err
	}
	outClean := strings.Trim(outString, "\n")
	dependenciesNames := strings.Split(outClean, "\n")
	if len(dependenciesNames) > 0 {
		for _, dep := range dependenciesNames {
			if dep == "" {
				continue
			}
			depSpell, err := m.bareRunSingle(cranSpell{
				PackageName:  dep,
				Repository:   c.Repository,
				BioConductor: false,
			}, bin, c.BioConductor)
			if err != nil {
				return nil, err
			}
			if !Contains(finalSpell.Dependencies, *depSpell) {
				finalSpell.Dependencies = append(finalSpell.Dependencies,
					*depSpell)
			}
		}
	}
	return finalSpell, nil
}

func (m *cranModule) bareRunSingle(
	c cranSpell,
	bin string,
	bioconductor bool,
) (*cranSpell, error) {
	tempdir := os.TempDir()
	downloadOptions := &rcran.DownloadOptions{
		PackageName:          c.PackageName,
		DestinationDirectory: tempdir,
		Repository:           c.Repository,
	}
	var cmd string
	var err error
	if bioconductor {
		cmd, err = rcran.DownloadBioconductor(downloadOptions)
	} else {
		cmd, err = rcran.Download(downloadOptions)
	}
	if err != nil {
		return nil, err
	}
	script, err := rscript.New(bin).Evaluate(cmd).Seal()
	if err != nil {
		return nil, err
	}
	// TODO: Check around here.
	out, err := script.CombinedOutput()
	logger.Get().Print(string(out))
	if err != nil {
		return nil, err
	}
	finalSpell := m.spellFromDownloadOptions(downloadOptions)
	finalSpell.BioConductor = bioconductor
	path, err := rcran.ParsePath(string(out))
	if err != nil {
		return nil, err
	}
	finalSpell.PackagePath = path
	return finalSpell, nil
}

// installBioconductor runs the command in an opinionated fashion to install
// bioconductor.
//
// Used if cranSpell.BioConductor is set.
func (m *cranModule) installBioconductor(cfg *Config) error {
	spell, err := m.bareRun(cranSpell{
		PackageName:  "BiocManager",
		BioConductor: false,
	}, cfg)
	if err != nil {
		return err
	}
	err = m.Commit(cfg, spell)
	if err == ErrAlreadyPresent {
		return nil
	}
	return err
}

// cobraArgs is used to validate the arguments passed to the cran command.
//
// This function is intended to be used by cobra.
func (c *cranModule) cobraArgs() func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s module requires at least one argument.",
				cranModuleName)
		}
		return nil
	}
}

// cobraRun is used to run the module from the command line.
// It serves as an entry point to the cranModule.
//
// This function is inteded to be used by cobra.
func (c *cranModule) cobraRun(cfg *Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		isBioconductor := strings.Compare(
			cmd.CalledAs(),
			bioconductorModuleName) == 0
		packageName := args[0]
		repository, _ := cmd.Flags().GetString("repository")
		spell, err := c.bareRun(cranSpell{
			PackageName:  packageName,
			Repository:   repository,
			BioConductor: isBioconductor,
		}, cfg)
		if err != nil {
			logger.Get().Print(err)
			return
		}
		err = c.Commit(cfg, spell)
		if err != nil {
			logger.Get().Print(err)
		}
	}
}

// CliConfig implements Module.
func (c *cranModule) CliConfig(config *Config) *cobra.Command {
	command := &cobra.Command{
		Args:    c.cobraArgs(),
		Example: cranModuleExample,
		Run:     c.cobraRun(config),
		Short:   cranModuleShort,
		Use:     cranModuleName,
		Aliases: []string{bioconductorModuleName},
	}
	command.PersistentFlags().String("repository", "", "Repository to use.")
	return command
}

// Commit implements Module.
func (c *cranModule) Commit(config *Config, result any) error {
	newEntry, ok := result.(*cranSpell)
	if !ok {
		return fmt.Errorf("Error Converting") //TODO: unify errors.
	}
	if Contains(config.Cran, *newEntry) {
		return ErrAlreadyPresent
	}
	config.Cran = append(config.Cran, *newEntry)
	return nil
}

// Run implements Module.
func (c *cranModule) Run(anyspell any) error {
	spell, ok := anyspell.(cranSpell)
	if !ok {
		return fmt.Errorf("Error converting")
	}
	project, err := project.ProjectPath()
	if err != nil {
		return err
	}
	bin, err := rscript.DetectRscriptBinary()
	if err != nil {
		return err
	}
	destinationDirectory := path.Join(*project, cranModuleName)
	err = os.MkdirAll(destinationDirectory, 0755)
	if err != nil {
		return err
	}
	downloadOptions := &rcran.DownloadOptions{
		PackageName:          spell.PackageName,
		DestinationDirectory: destinationDirectory,
		Repository:           spell.Repository,
	}
	var cmd string
	if spell.BioConductor {
		cmd, err = rcran.DownloadBioconductor(downloadOptions)
	} else {
		cmd, err = rcran.Download(downloadOptions)
	}
	if err != nil {
		return err
	}
	script, err := rscript.New(bin).Evaluate(cmd).Seal()
	script.Stdout = os.Stdout
	script.Stderr = os.Stderr
	err = script.Run()
	return err
}
