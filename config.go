package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Shell string
	Files []File
}

type File struct {
	Path   string
	Name   string
	Create string
	Change string
	Delete string
}

func ReadConfig(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = yaml.Unmarshal(file, config)

	return config, err
}
