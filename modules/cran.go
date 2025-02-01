package modules

import (
	"bytes"
	"credo/cache"
	"credo/logger"
	"credo/project"
	"credo/suggest"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	gorcran "github.com/CREDOProject/go-rcran"
	gordepends "github.com/CREDOProject/go-rdepends"
	gordependsP "github.com/CREDOProject/go-rdepends/providers"
	gorscript "github.com/CREDOProject/go-rscript"
	"github.com/CREDOProject/sharedutils/filter"
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

type cranSpell struct {
	PackageName          string      `yaml:"package_name,omitempty"`
	PackagePath          string      `yaml:"package_path,omitempty"`
	Repository           string      `yaml:"repository,omitempty"`
	BioConductor         bool        `yaml:"bioconductor,omitempty"`
	Dependencies         []cranSpell `yaml:"dependencies,omitempty"`
	ExternalDependencies Config      `yaml:"external_dependencies,omitempty"`
}

// Registers the cranModule.
func init() { Register(cranModuleName, func() Module { return &cranModule{} }) }

type cranModule struct{}

// Apply implements Module.
func (c *cranModule) Apply(anyspell any) error {
	var spell cranSpell
	switch s := anyspell.(type) {
	case cranSpell:
		spell = s
	case *cranSpell:
		if s == nil {
			return fmt.Errorf("[cran/apply]: %v", ErrConverting)
		}
		spell = *s
	default:
		return fmt.Errorf("[cran/apply]: %v", ErrConverting)
	}
	if cache.Retrieve(cranModuleName, spell.PackageName) != nil {
		return nil
	}
	err := DeepApply(&spell.ExternalDependencies)
	if err != nil {
		return err
	}
	for _, dep := range slices.Backward(spell.Dependencies) {
		err := c.Apply(&dep)
		if err != nil {
			return err
		}
	}
	destdir, err := c.destinationDirectory()
	if err != nil {
		return err
	}
	libraryDir, err := c.libraryDirectory()
	if err != nil {
		return err
	}
	localInstallOptions := gorcran.InstallOptions{
		PackageName: path.Join(destdir, spell.PackagePath),
		Repository:  "NULL",
		Library:     libraryDir,
	}
	cmd, err := gorcran.InstallLocal(&localInstallOptions)
	if err != nil {
		return err
	}
	bin, err := gorscript.DetectRscriptBinary()
	if err != nil {
		return err
	}
	script, err := gorscript.New(bin).Evaluate(cmd).Seal()
	script.Stdout = os.Stdout
	script.Stderr = os.Stderr
	err = script.Run()
	if err == nil {
		_ = cache.Insert(cranModuleName, spell.PackageName, true)
	}
	return err
}

// BulkApply implements Module.
func (c *cranModule) BulkApply(config *Config) error {
	for _, cs := range config.Cran {
		if err := c.Apply(cs); err != nil {
			return err
		}
	}
	return nil
}

// BulkSave implements Module.
func (c *cranModule) BulkSave(config *Config) error {
	for _, cs := range config.Cran {
		if err := c.Save(cs); err != nil {
			return err
		}
	}
	return nil
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
		return fmt.Errorf("[cran/commit]: %v", ErrConverting)
	}
	if newEntry == nil {
		return nil
	}
	if Contains(config.Cran, *newEntry) {
		return ErrAlreadyPresent
	}
	config.Cran = append(config.Cran, *newEntry)
	return nil
}

// Save implements Module.
func (c *cranModule) Save(anyspell any) error {
	spell, ok := anyspell.(cranSpell)
	if !ok {
		return fmt.Errorf("[cran/save]: %v", ErrConverting)
	}
	for _, dep := range spell.Dependencies {
		if err := c.Save(dep); err != nil {
			return err
		}
	}
	err := DeepSave(&spell.ExternalDependencies)
	if err != nil {
		return nil
	}
	if cache.Retrieve(cranModuleName, spell.PackageName) != nil {
		return nil
	}
	destdir, err := c.destinationDirectory()
	if err != nil {
		return fmt.Errorf(`[cran] dest: %v`, err)
	}
	filesMap, err := listDownloadedFilesInMap(destdir)
	if err != nil {
		return fmt.Errorf(`[cran] list: %v`, err)
	}
	if _, present := filesMap[spell.PackagePath]; present {
		logger.Get().Printf(`[cran]: Skipped saving %s, already present.`,
			spell.PackageName)
		return nil
	}
	downloadFunction := c.downloadFunction(spell.BioConductor)
	libraryDir, err := c.libraryDirectory()
	if err != nil {
		return err
	}
	cmd, err := downloadFunction(&gorcran.DownloadOptions{
		Library:              libraryDir,
		PackageName:          spell.PackageName,
		DestinationDirectory: destdir,
		Repository:           spell.Repository,
	})
	if err != nil {
		return err
	}
	bin, err := gorscript.DetectRscriptBinary()
	if err != nil {
		return err
	}
	script, err := gorscript.New(bin).Evaluate(cmd).Seal()
	script.Stdout = os.Stdout
	script.Stderr = os.Stderr
	err = script.Run()
	if err == nil {
		_ = cache.Insert(cranModuleName, spell.PackageName, true)
	}
	return err
}

func (c *cranModule) destinationDirectory() (string, error) {
	project, err := project.ProjectPath()
	if err != nil {
		return "", err
	}
	directory := path.Join(*project, cranModuleName)
	err = os.MkdirAll(directory, 0755)
	if err != nil {
		return "", err
	}
	return directory, nil
}

func listDownloadedFilesInMap(destdir string) (map[string]struct{}, error) {
	filesMap := make(map[string]struct{})
	err := filepath.Walk(destdir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filesMap[info.Name()] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return filesMap, nil
}

func (c *cranModule) libraryDirectory() (string, error) {
	project, err := project.ProjectPath()
	if err != nil {
		return "", err
	}
	directory := path.Join(*project, "R-Library")
	err = os.MkdirAll(directory, 0755)
	if err != nil {
		return "", err
	}
	return directory, nil
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
// This function is intended to be used by cobra.
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

func (c *cranModule) bareRun(s cranSpell, cfg *Config) (*cranSpell, error) {
	if s.BioConductor {
		err := c.installBioConductor(cfg)
		if err != nil {
			return nil, err
		}
	}
	return c.bareRunSingle(s)
}

func (c *cranModule) bareRunSingle(s cranSpell) (*cranSpell, error) {
	if spell := cache.Retrieve(cranModuleName,
		s.PackageName); spell != nil {
		if newSpell, ok := spell.(cranSpell); ok {
			return &newSpell, nil
		}
	}
	rscriptBin, err := gorscript.DetectRscriptBinary()
	if err != nil {
		return nil, fmt.Errorf(`[cran] detect: %v`, err)
	}
	tempdir := os.TempDir()
	libraryDir, err := c.libraryDirectory()
	if err != nil {
		return nil, err
	}
	downloadOptions := &gorcran.DownloadOptions{
		Library:              libraryDir,
		PackageName:          s.PackageName,
		DestinationDirectory: tempdir,
		Repository:           s.Repository,
	}
	downloadFunction := c.downloadFunction(s.BioConductor)
	cmd, err := downloadFunction(downloadOptions)
	if err != nil {
		return nil, fmt.Errorf(`[cran] download: %v`, err)
	}
	script, err := gorscript.New(rscriptBin).Evaluate(cmd).Seal()
	if err != nil {
		return nil, fmt.Errorf(`[cran] generate: %v`, err)
	}
	var buffer bytes.Buffer
	script.Stdout = io.MultiWriter(os.Stdout, &buffer)
	script.Stderr = os.Stderr
	err = script.Run()
	if err != nil {
		return nil, fmt.Errorf(`[cran] run: %v`, err)
	}
	outputString := buffer.String()
	pkgPath, err := gorcran.GetPath(outputString)
	if err != nil {
		return nil, fmt.Errorf(`[cran] path: %v`, err)
	}
	deps, err := c.getDependencies(rscriptBin, s)
	if err != nil {
		return nil, fmt.Errorf(`[cran] deps: %v`, err)
	}
	additionalDependencies, err := gordepends.DependsOn(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("[cran] rdepends: %v", err)
	}
	// Register suggestions.
	suggestions := filter.Filter(additionalDependencies,
		func(a gordependsP.Dependency) bool { return a.Suggestion })
	for _, suggestion := range suggestions {
		suggest.Register(suggest.Suggestion{
			Module:    cranModuleName,
			From:      s.PackageName,
			Suggested: suggestion.Name,
		})
	}
	parsedPath, err := gorcran.ParsePath(outputString)
	if err != nil {
		return nil, fmt.Errorf(`[cran] parse: %v`, err)
	}
	finalSpell := cranSpell{
		PackageName:  s.PackageName,
		PackagePath:  parsedPath,
		Repository:   s.Repository,
		BioConductor: s.BioConductor,
		Dependencies: deps,
	}
	for _, d := range additionalDependencies {
		module, ok := Modules[d.PackageManager]
		if ok {
			args := []string{d.Name}
			module().CliConfig(&finalSpell.ExternalDependencies).Run(nil, args)
		}
	}
	_ = cache.Insert(cranModuleName, s.PackageName, finalSpell)
	return &finalSpell, nil
}

func (c *cranModule) getDependencies(rscriptBin string, s cranSpell) ([]cranSpell, error) {
	dependencyFunction := c.dependencyFunction(s.BioConductor)
	libraryDir, err := c.libraryDirectory()
	if err != nil {
		return nil, err
	}
	cmd, err := dependencyFunction(&gorcran.InstallOptions{
		PackageName: s.PackageName,
		Repository:  s.Repository,
		DryRun:      false,
		Library:     libraryDir,
	})
	if err != nil {
		return nil, err
	}
	script, err := gorscript.New(rscriptBin).Evaluate(cmd).Seal()
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	script.Stdout = &buffer
	script.Stderr = os.Stderr
	err = script.Run()
	if err != nil {
		return nil, err
	}
	outputString := buffer.String()
	dependencyList := strings.Split(strings.Trim(outputString, "\n"), "\n")
	deps := []cranSpell{}
	var MaxWorkers chan int = make(chan int, 4)
	var wg sync.WaitGroup
	for _, dep := range dependencyList {
		if dep == "" {
			continue
		}
		wg.Add(1)
		MaxWorkers <- 1
		go func(dep string) {
			defer func() { wg.Done(); <-MaxWorkers }()
			fmt.Printf("Worker %s starting\n", dep)
			dependencySpell, err := c.bareRunSingle(cranSpell{
				PackageName:  dep,
				Repository:   s.Repository,
				BioConductor: s.BioConductor,
			})
			if err != nil || dependencySpell == nil {
				return
			}
			if !Contains(deps, *dependencySpell) {
				deps = append(deps, *dependencySpell)
			}
		}(dep)
	}
	wg.Wait()
	return deps, nil
}

func (c *cranModule) dependencyFunction(bioconductor bool) func(
	o *gorcran.InstallOptions) (string, error) {
	if bioconductor {
		return gorcran.GetBioconductorDependencies
	}
	return gorcran.GetDependencies
}

func (c *cranModule) downloadFunction(bioconductor bool) func(
	o *gorcran.DownloadOptions) (string, error) {
	if bioconductor {
		return gorcran.DownloadBioconductor
	}
	return gorcran.Download
}

func (c *cranModule) installBioConductor(cfg *Config) error {
	spell, err := c.bareRun(cranSpell{
		PackageName:  "BiocManager",
		BioConductor: false,
	}, cfg)
	if err != nil {
		return err
	}
	err = c.Commit(cfg, spell)
	if err != ErrAlreadyPresent && err != nil {
		return err
	}
	err = c.Save(*spell)
	if err != nil {
		return err
	}
	err = c.Apply(*spell)
	return err
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
