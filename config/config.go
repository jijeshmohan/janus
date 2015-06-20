package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/jijeshmohan/janus/rest"
)

// Config type represent configuration json.
type Config struct {
	Port      int             `json:"port,omitempty"`
	Delay     int             `json:"delay,omitempty"`
	Auth      *auth           `json:"auth,omitempty"`
	Path      string          `json:"-"`
	Resources []rest.Resource `json:"resources,omitempty"`
	URLs      []rest.URL      `json:"urls,omitempty"`
	EnableLog bool            `json:"enableLog,omitempty"`
}

type auth struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

// ParseFile parse input file and generate Config type.
// returns error incase of invalid file or format.
func ParseFile(path string) (*Config, error) {

	// Read file
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("Unable to open config.json. Please check the file is present")
	}

	defer file.Close()

	config := Config{}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		return nil, err
	}
	return &config, nil
}
