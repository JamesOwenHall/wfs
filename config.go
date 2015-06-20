package main

import (
	"errors"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Shell string
	Paths map[string]map[string]string
}

func ReadConfig(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	if err = yaml.Unmarshal(file, config); err != nil {
		return nil, err
	}

	normalize(config)

	if err = validate(config); err != nil {
		return nil, err
	}

	return config, nil
}

// normalize transforms all of the actions in the Paths map to have keys which
// are all lowercase.
func normalize(config *Config) {
	newMap := make(map[string]map[string]string)

	for path, m := range config.Paths {
		newMap[path] = make(map[string]string)

		for event, action := range m {
			newMap[path][strings.ToLower(event)] = action
		}
	}

	config.Paths = newMap
}

// validate returns an error if any event is not either "create", "change" or
// "delete".
func validate(config *Config) error {
	for path, m := range config.Paths {
		for event, _ := range m {
			if event != "create" && event != "change" && event != "delete" {
				return errors.New("unknown event " + event + " in path " + path)
			}
		}
	}

	return nil
}
