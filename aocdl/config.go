package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type configuration struct {
	SessionCookie      string `json:"session-cookie"`
	Output             string `json:"output"`
	Year               int    `json:"year"`
	Day                int    `json:"day"`
	Force              bool   `json:"-"`
	Wait               bool   `json:"-"`
	StoryOut           string `json:"story-out"`
	Template           string `json:"template"`
	TemplateOutput     string `json:"template-output"`
	TestOutput         string `json:"test-output"`
	TestTemplate       string `json:"test-template"`
	TestTemplateOutput string `json:"test-template-output"`
}

func loadConfigs() (*configuration, error) {
	config := new(configuration)

	home := ""
	usr, err := user.Current()
	if err == nil {
		home = usr.HomeDir
	}

	if home != "" {
		err = config.mergeWithFileIfExists(filepath.Join(home, ".aocdlconfig"))
		if err != nil {
			return nil, err
		}
	}

	wd, _ := os.Getwd()

	// If we could not determine either directory or if we are not currently in
	// the home directory, try and load the configuration relative to the
	// current working directory.
	if wd == "" || home == "" || wd != home {
		err = config.mergeWithFileIfExists(".aocdlconfig")
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func loadConfig(filename string) (*configuration, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := new(configuration)
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (config *configuration) mergeWithFileIfExists(filename string) error {
	loaded, err := loadConfig(filename)
	if err == nil {
		// file loaded
		config.merge(loaded)
		return nil
	} else if os.IsNotExist(err) {
		// file not found
		return nil
	} else {
		// read error
		return err
	}
}

func (config *configuration) merge(other *configuration) {
	if other.Force {
		config.Force = true
	}
	if other.SessionCookie != "" {
		config.SessionCookie = other.SessionCookie
	}
	if other.Output != "" {
		config.Output = other.Output
	}
	if other.Year != 0 {
		config.Year = other.Year
	}
	if other.Day != 0 {
		config.Day = other.Day
	}
	if other.StoryOut != "" {
		config.StoryOut = other.StoryOut
	}
	if other.TestOutput != "" {
		config.TestOutput = other.TestOutput
	}
	if other.TestTemplate != "" {
		config.TestTemplate = other.TestTemplate
	}
	if other.TestTemplateOutput != "" {
		config.TestTemplateOutput = other.TestTemplateOutput
	}
	if other.TestTemplateOutput != "" {
		config.TestTemplateOutput = other.TestTemplateOutput
	}
	if other.TemplateOutput != "" {
		config.TemplateOutput = other.TemplateOutput
	}
	if other.Template != "" {
		config.Template = other.Template
	}
}
