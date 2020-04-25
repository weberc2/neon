package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

type Config struct {
	InputDirectory      string      `json:"input_directory"`
	LiveReloadPort      uint16      `json:"live_reload_port"`
	OutputDirectory     string      `json:"output_directory"`
	OutputExtension     string      `json:"output_exptension"`
	PostOutputDirectory string      `json:"post_output_directory"`
	SiteRoot            string      `json:"site_root"`
	Theme               string      `json:"theme"`
	ThemeConfig         interface{} `json:"theme_config"`

	// See https://github.com/alecthomas/chroma/tree/master/styles for complete
	// list.
	CodeHighlightTheme string `json:"code_highlight_theme"`
	IndexPageSize      int    `json:"index_page_size"`
}

func loadConfig(filepath string, conf *Config) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, conf)
}

func findRoot(from, match string) (string, error) {
	if from == "/" {
		return "", nil
	}

	files, err := ioutil.ReadDir(from)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.Name() == match {
			return from, nil
		}
	}

	return findRoot(filepath.Dir(from), match)
}

func Load() (Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err) // not much we can do about this
	}

	root, err := findRoot(wd, "neon.yaml")
	if err != nil {
		// If we get an error back, it must mean there was something wrong with
		// the file system, or there is an error in the program's logic that
		// allowed us to pass a non-existant or otherwise bad path to
		// findRoot().
		panic(err)
	}
	if root == "" {
		return Config{}, errors.New(
			"Couldn't find a parent directory with a neon.yaml file.",
		)
	}
	conf := Config{
		InputDirectory:      root,
		OutputDirectory:     filepath.Join(root, "_output"),
		PostOutputDirectory: "posts",
		Theme:               "default",
		OutputExtension:     "html",
		SiteRoot:            "http://localhost:8000/",
		IndexPageSize:       3,
	}

	// TODO: config validation, handle required fields, etc
	err = loadConfig(filepath.Join(root, "neon.yaml"), &conf)
	return conf, err
}
