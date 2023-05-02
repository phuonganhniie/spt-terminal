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
	const filePerm = 0755
	if _, err := os.Stat(app.userConfigPath); os.IsNotExist(err) {
		err := os.MkdirAll(app.userConfigPath, filePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
