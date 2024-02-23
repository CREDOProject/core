package project

import (
	"io/fs"
	"os"
	"path"
)

// Global project path.
var gPath *string

// Gets the project path.
func ProjectPath() (*string, error) {
	if gPath != nil {
		return gPath, nil
	}
	basePath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	projectPath := path.Join(basePath, "credoenv")

	// Create project path.
	err = os.Mkdir(projectPath, fs.ModeDir)
	if err != nil {
		return nil, err
	}
	gPath = &projectPath

	return &projectPath, nil
}
