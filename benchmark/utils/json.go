package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ReadJSON will read and serialize a json file into a given target struct.
func ReadJSON(path string, target interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

// SaveAsJSON saves a struct as a json file to the given path.
// If the path doesn't exists creates it.
func SaveAsJSON(path string, source interface{}) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	asJSON, err := json.MarshalIndent(source, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, asJSON, 0644)
}
