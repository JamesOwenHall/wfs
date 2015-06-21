package main

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

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

// ReadConfig reads the configuration file and returns the *Config that
// represents it, or an error if there's a problem parsing the file.
func ReadConfig(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = yaml.Unmarshal(file, config)

	// We only deal with clean paths.
	for _, file := range config.Files {
		file.Path = filepath.Clean(file.Path)
	}

	return config, err
}
