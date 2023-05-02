package config

import (
	"os"
	"path/filepath"
)

// UserConfigPath returns the path to the configuration file.
func userConfigPath() string {
	userConfigDir, _ := userConfigDir()
	return filepath.Join(userConfigDir, "stp-terminal")
}

// UserConfigDir returns the path to the configuration directory.
func userConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	return configDir, err
}

// UserCacheDir returns the path to the cache directory.
func userCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	return cacheDir, err
}

func (app *AppConfig) createUserConfigPath() error {
	if _, err := os.Stat(app.UserConfigPath); os.IsNotExist(err) {
		err := os.MkdirAll(app.UserConfigPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
