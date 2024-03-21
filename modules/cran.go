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
	credo bioconductor <>
`

// Registers the carnModule.
func init() { Register(cranModuleName, func() Module { return &cranModule{} }) }

// cranModule is used to manage the CARN scope in the credospell configuration.
type cranModule struct{}

type cranSpell struct {
	PackageName  string `yaml:"package_name,omitempty"`
	PackagePath  string `yaml:"package_path,omitempty"`
	Repository   string `yaml:"repository,omitempty"`
	BioConductor bool   `yaml:"bioconductor,omitempty"`
}

// spellFromInstallOptions returns a new cran spell from *rcran.InstallOptions.
func (*cranModule) spellFromInstallOptions(
	options *rcran.InstallOptions,
) *cranSpell {
	spell := &cranSpell{
		PackageName: options.PackageName,
		Repository:  options.Repository,
	}
	return spell
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
	return strings.Compare(s.PackageName, c.PackageName) == 0 &&
		strings.Compare(s.PackagePath, c.PackagePath) == 0
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
	tempdir := os.TempDir()
	downloadOptions := &rcran.DownloadOptions{
		PackageName:          c.PackageName,
		DestinationDirectory: tempdir,
		Repository:           c.Repository,
	}
	cmd, err := rcran.Download(downloadOptions)
	if err != nil {
		return nil, err
	}
	script, err := rscript.New(bin).Evaluate(cmd).Seal()
	if err != nil {
		return nil, err
	}
	out, err := script.CombinedOutput()
	if err != nil {
		return nil, err
	}
	logger.Get().Print(string(out))
	finalSpell := m.spellFromDownloadOptions(downloadOptions)
	path, err := rcran.ParsePath(string(out))
	if err != nil {
		return nil, err
	}
	if path == "" {
		return nil, fmt.Errorf("Error downloading package.")
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
			logger.Get().Fatal(err)
		}
		err = c.Commit(cfg, spell)
		if err != nil {
			logger.Get().Fatal(err)
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
	cmd, err := rcran.Download(downloadOptions)
	if err != nil {
		return err
	}
	script, err := rscript.New(bin).Evaluate(cmd).Seal()
	script.Stdout = os.Stdout
	script.Stderr = os.Stderr
	err = script.Run()
	return err
}
