package config

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

func Path() (string, error) {
	xdgPath := os.Getenv("XDG_CONFIG_HOME")
	configPath := xdgPath
	if configPath == "" {
		var err error
		configPath, err = homedir.Dir()
		if err != nil {
			return "", err
		}
	}
	return filepath.Join(configPath, ".svix"), nil
}

func Write(settings interface{}) error {
	// write config file manually because viper cant write to dotfiles without extensions
	// see: https://github.com/spf13/viper/pull/1064
	cfgPath, err := Path()
	if err != nil {
		return err
	}
	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	f, err := os.OpenFile(cfgPath, flags, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}
	if _, err = f.Write(b); err != nil {
		return err
	}
	return nil
}
