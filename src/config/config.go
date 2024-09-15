package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/thehxdev/aptgit/gpath"
)

var (
	G *Config = &Config{}

	lookupPaths []string = []string{
		"./config.json",
		"~/.aptgit/config.json",
		"/etc/aptgit/config.json",
	}
)

type Config struct {
	Home string `json:"home,omitempty"`

	// These other paths MUST be relative to `Home` field
	BinPath      string
	InstallPath  string
	DownloadPath string
	Gpkgs        string
}

func ReadConfig(p string) error {
	var configPath string = p
	var err error = nil

	if configPath == "" {
		configPath, err = findConfigFile()
		if err != nil {
			return err
		}
	}

	fp, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer fp.Close()

	err = json.NewDecoder(fp).Decode(G)
	if err != nil {
		return err
	}

	if G.Home == "" {
		return errors.New("home field is empty in config file")
	}

	G.Home = gpath.Expand(G.Home)
	G.InstallPath = path.Join(G.Home, "installs")
	G.DownloadPath = path.Join(G.Home, "downloads")
	G.BinPath = path.Join(G.Home, "bin")
	G.Gpkgs = path.Join(G.Home, "gpkgs")

	return nil
}

func findConfigFile() (string, error) {
	for _, p := range lookupPaths {
		_, err := os.Stat(gpath.Expand(p))
		if err == nil {
			return p, nil
		}
	}
	return "", errors.New("config file not found")
}
