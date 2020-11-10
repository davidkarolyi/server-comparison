package report

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
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
