package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	ErrNoShell = fmt.Errorf("no shell defined in configuration.")
	ErrNoFiles = fmt.Errorf("no file patterns defined in configuration.")
)

type ErrBadFile error

// Config is an internal representation of a configuration file.  Shell is
// the application that gets run when an event occurs.
type Config struct {
	Shell string
	Files []File
}

// File is an internal representation of a category of files on disk to watch.
// Path is the directory to watch, while Name is the search pattern for files
// within the Path.  Create, Change and Delete are the commands that are run
// whenever a file that matches the Path and Name is created, changed or
// deleted, respectively.
type File struct {
	Path   string
	Name   string
	Create string
	Change string
	Delete string
}

// ReadConfigFile reads the configuration file and returns the *Config that
// represents it, or an error if there's a problem reading or parsing the file.
func ReadConfigFile(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ReadConfig(file)
}

// ReadConfig raeds the configuration data and returns the *Config that
// represents it, or an error if there's a problem parsing the data.
func ReadConfig(input []byte) (*Config, error) {
	config := new(Config)
	if err := yaml.Unmarshal(input, config); err != nil {
		return nil, err
	}

	trimSpace(config)
	if err := validate(config); err != nil {
		return nil, err
	}

	// We only deal with clean paths.
	for i := range config.Files {
		file := &config.Files[i]
		file.Path = filepath.Clean(file.Path)
	}

	return config, nil
}

// trimSpace removes leading and trailing white space from all fields of
// config.
func trimSpace(config *Config) {
	config.Shell = strings.TrimSpace(config.Shell)
	for i := range config.Files {
		file := &config.Files[i]
		file.Path = strings.TrimSpace(file.Path)
		file.Name = strings.TrimSpace(file.Name)
		file.Create = strings.TrimSpace(file.Create)
		file.Change = strings.TrimSpace(file.Change)
		file.Delete = strings.TrimSpace(file.Delete)
	}
}

// validate returns an error if config is invalid.  It returns either
// ErrNoShell, ErrNoFiles or ErrBadFile.
func validate(config *Config) error {
	if config.Shell == "" {
		return ErrNoShell
	}

	if len(config.Files) == 0 {
		return ErrNoFiles
	}

	for _, file := range config.Files {
		if file.Path == "" {
			return ErrBadFile(fmt.Errorf("no path defined for a file pattern."))
		}

		if file.Name == "" {
			return ErrBadFile(fmt.Errorf(
				"no name pattern defined for file at path %s",
				file.Path,
			))
		}

		if file.Create == "" && file.Change == "" && file.Delete == "" {
			return ErrBadFile(fmt.Errorf(
				"no actions defined for file at path %v",
				file.Path,
			))
		}
	}

	return nil
}
