package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SaveAsJSON saves a struct as a json file to the given path.
// If the path doesn't exists creates it.
func SaveAsJSON(path string, v interface{}) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	asJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, asJSON, 0644)
}

// ChangeToProjectRoot will set the working directory to the project's root if it's exists.
func ChangeToProjectRoot() error {
	rootPath, err := ProjectRoot()
	if err != nil {
		return err
	}

	return os.Chdir(rootPath)
}

// ProjectRoot will return the path to the project's root folder
func ProjectRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	path := new(bytes.Buffer)
	cmd.Stdout = path

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(path.String()), nil
}
