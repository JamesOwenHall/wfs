package main

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	var (
		input    []byte
		expected *Config
		actual   *Config
		err      error
	)

	// Empty input
	input = []byte(``)
	_, err = ReadConfig(input)
	if err != ErrNoShell && err != ErrNoFiles {
		t.Errorf("Expected error, got %v", err)
	}

	// Empty fields
	input = []byte(`
shell: 
files: `)
	_, err = ReadConfig(input)
	if err != ErrNoShell && err != ErrNoFiles {
		t.Errorf("Expected error, got %v", err)
	}

	// Missing shell
	input = []byte(`
files: 
  - path: .`)
	_, err = ReadConfig(input)
	if err != ErrNoShell {
		t.Errorf("Expected ErrNoShell, got %v", err)
	}

	// Missing files
	input = []byte(`
shell: bash`)
	_, err = ReadConfig(input)
	if err != ErrNoFiles {
		t.Errorf("Expected ErrNoFiles, got %v", err)
	}

	// Missing file path
	input = []byte(`
shell: bash
files: 
  - name: "*"
    create: echo hello`)
	actual, err = ReadConfig(input)
	if _, ok := err.(ErrBadFile); !ok {
		t.Log(actual)
		t.Errorf("Expected ErrBadFile, got %v", err)
	}

	// Missing file name
	input = []byte(`
shell: bash
files: 
  - path: .
    create: echo hello`)
	_, err = ReadConfig(input)
	if _, ok := err.(ErrBadFile); !ok {
		t.Errorf("Expected ErrBadFile, got %v", err)
	}

	// Missing actions
	input = []byte(`
shell: bash
files: 
  - path: .
    name: "*"`)
	_, err = ReadConfig(input)
	if _, ok := err.(ErrBadFile); !ok {
		t.Errorf("Expected ErrBadFile, got %v", err)
	}

	// Basic config
	input = []byte(`
shell: bash
files: 
  - path: .
    name: "*"
    create: echo $file has been created >> out.log`)
	expected = &Config{
		Shell: "bash",
		Files: []File{
			File{Path: ".", Name: "*", Create: "echo $file has been created >> out.log"},
		},
	}
	actual, err = ReadConfig(input)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !equalConfigs(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	// Multiline action
	input = []byte(`
shell: bash
files: 
  - path: .
    name: "*"
    create: |
      echo $file has been created >> out.log
      echo second command >> out.log`)
	expected = &Config{
		Shell: "bash",
		Files: []File{
			File{Path: ".", Name: "*", Create: "echo $file has been created >> out.log\necho second command >> out.log"},
		},
	}
	actual, err = ReadConfig(input)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !equalConfigs(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	// Multiple files
	input = []byte(`
shell: bash
files: 
  - path: .
    name: "*"
    create: echo $file has been created >> out.log
  - path: ./sub
    name: "*.txt"
    change: echo $file has been changed >> out.log`)
	expected = &Config{
		Shell: "bash",
		Files: []File{
			File{Path: ".", Name: "*", Create: "echo $file has been created >> out.log"},
			File{Path: "sub", Name: "*.txt", Change: "echo $file has been changed >> out.log"},
		},
	}
	actual, err = ReadConfig(input)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !equalConfigs(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func equalConfigs(a, b *Config) bool {
	if a == nil && b == nil {
		return true
	}

	// We know they're not both nil, so if one of them is, then they're not
	// equal.
	if a == nil || b == nil {
		return false
	}

	if a.Shell != b.Shell {
		return false
	}

	if len(a.Files) != len(b.Files) {
		return false
	}

	for i := range a.Files {
		aFile, bFile := a.Files[i], b.Files[i]
		if !(aFile.Path == bFile.Path && aFile.Name == bFile.Name &&
			aFile.Create == bFile.Create && aFile.Change == bFile.Change &&
			aFile.Delete == bFile.Delete) {
			return false
		}
	}

	return true
}
