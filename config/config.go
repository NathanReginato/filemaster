package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

const configFilePath = "./config.yaml"

// Config represents the YAML file that configures this app
type Config struct {
	Workspace string            `yaml:"workspace"`
	Structure []string          `yaml:"file-structure"`
	Process   []string          `yaml:"process"`
	Alias     map[string]string `yaml:"aliases"`
}

// New return a Config struct populated with the values from the config file
func New() (*Config, error) {
	c := Config{}

	absPath, _ := filepath.Abs(configFilePath)

	dat, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %v", err)
	}

	err = yaml.Unmarshal([]byte(dat), &c)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal config file: %v", err)
	}

	return &c, nil
}

// GetWorkspace return the workspace directory
func (c *Config) GetWorkspace() string {
	// TODO: Account for unintialized File
	return c.Workspace
}

// GetStructure return the directory structure as layed out in the configuration
func (c *Config) GetStructure() []string {
	// TODO: Account for unintialized File
	return c.Structure
}

// HasDir will iterate over all the keys in the sturcture array and return true if the passed in value is found.
func (c *Config) HasDir(e string) bool {
	for _, v := range c.Structure {
		if v == e {
			return true
		}
	}
	return false
}

func (c *Config) GetAlias() map[string]string {
	return c.Alias
}
