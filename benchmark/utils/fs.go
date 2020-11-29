package utils

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

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

// ListDirContent returns a list of files and a list of directory names
// located in the given path.
func ListDirContent(path string) (fileNames []string, dirNames []string, err error) {
	filesAndDirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}

	for _, fileOrDir := range filesAndDirs {
		if fileOrDir.IsDir() {
			dirNames = append(dirNames, fileOrDir.Name())
		} else {
			fileNames = append(fileNames, fileOrDir.Name())
		}
	}

	return fileNames, dirNames, nil
}
