package utils

import (
	"encoding/json"
	"os"
)

// SaveToJSON guarda datos en un archivo JSON
func SaveToJSON(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// LoadFromJSON carga datos desde un archivo JSON
func LoadFromJSON(filename string, target interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(target)
}