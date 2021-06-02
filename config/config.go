package config

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

const FileName = "config.toml"
const FileMode = os.FileMode(0600)

func Folder() (string, error) {
	var configPath string
	switch runtime.GOOS {
	case "windows":
		configPath = os.Getenv("APPDATA")
	default:
		configPath = os.Getenv("XDG_CONFIG_HOME")
	}
	if configPath == "" {
		var err error
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		configPath = filepath.Join(home, ".config")
	}
	return filepath.Join(configPath, "svix"), nil
}

func Write(settings map[string]interface{}) error {
	// write config file manually because viper cant write to dotfiles without extensions
	// see: https://github.com/spf13/viper/pull/1064
	cfgPath, err := Folder()
	if err != nil {
		return err
	}
	err = os.MkdirAll(cfgPath, os.FileMode(0700))
	if err != nil {
		return err
	}

	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	f, err := os.OpenFile(filepath.Join(cfgPath, FileName), flags, FileMode)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(settings); err != nil {
		return err
	}
	if _, err = buf.WriteTo(f); err != nil {
		return err
	}
	return nil
}
