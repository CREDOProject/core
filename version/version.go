package version

import (
	"os"
	"text/template"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

const versionTemplate = `{{.Name}}:
    Version:    {{.Version}}
    Commit:     {{.Commit}}
    Build Date: {{.BuildDate}}
`

type versionInfo struct {
	Name      string
	Version   string
	Commit    string
	BuildDate string
}

func PrintVersion(name string) error {
	info := versionInfo{
		Name:      name,
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
	}

	tmpl, err := template.New("version").Parse(versionTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stdout, info)
}
