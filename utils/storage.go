package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const DataDir = "data"

func ensureDataDir() error {
	return os.MkdirAll(DataDir, 0755)
}

func LoadJSON(filename string, v interface{}) error {
	if err := ensureDataDir(); err != nil {
		return err
	}
	path := filepath.Join(DataDir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, v)
}

func SaveJSON(filename string, v interface{}) error {
	if err := ensureDataDir(); err != nil {
		return err
	}
	path := filepath.Join(DataDir, filename)
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
