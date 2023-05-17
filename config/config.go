package config

import (
	"fmt"
	"path/filepath"

	"github.com/aditya-K2/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const PROJECT_SUB_DIR = "stp-terminal/"

type AppConfig struct {
	configDir      string
	UserCacheDir   string
	UserConfigPath string
	OnConfigChange func()
	config         *Config
	configErr      error
	cacheErr       error
}

var App = NewAppConfig()

func NewAppConfig() *AppConfig {
	configDir, configErr := userConfigDir()
	userCacheDir, cacheErr := userCacheDir()
	userConfigPath := filepath.Join(configDir, PROJECT_SUB_DIR)
	config := NewConfig()

	return &AppConfig{
		configDir:      configDir,
		UserCacheDir:   userCacheDir,
		UserConfigPath: userConfigPath,
		config:         config,
		configErr:      configErr,
		cacheErr:       cacheErr,
	}
}

func NewConfig() *Config {
	userCacheDir, _ := userCacheDir()
	return &Config{
		CacheDir:       utils.CheckDirectoryFmt(userCacheDir),
		RedrawInterval: 500,
		Colors:         NewColors(),
		HideImage:      false,
		RoundedCorners: false,
	}
}

type Config struct {
	CacheDir           string  `yaml:"cache_dir" mapstructure:"cache_dir"`
	RedrawInterval     int     `yaml:"redraw_interval" mapstructure:"redraw_interval"`
	Colors             *Colors `mapstructure:"colors"`
	AdditionalPaddingX int     `yaml:"additional_padding_x" mapstructure:"additional_padding_x"`
	AdditionalPaddingY int     `yaml:"additional_padding_y" mapstructure:"additional_padding_y"`
	ImageWidthExtraX   int     `yaml:"image_width_extra_x" mapstructure:"image_width_extra_x"`
	ImageWidthExtraY   int     `yaml:"image_width_extra_y" mapstructure:"image_width_extra_y"`
	HideImage          bool    `yaml:"hide_image" mapstructure:"hide_image"`
	RoundedCorners     bool    `yaml:"rounded_corners" mapstructure:"rounded_corners"`
}

func (app *AppConfig) ReadConfig() {
	app.parseFlags()
	app.checkConfigErrors()

	err := app.createUserConfigPath()
	if err != nil {
		utils.Print("RED", "Could not create user config path.\n")
		panic(err)
	}

	fileName := "config.yaml"
	viper.SetConfigFile(filepath.Join(app.UserConfigPath, fileName))

	if err := viper.ReadInConfig(); err != nil {
		errMsg := fmt.Sprintf("Could not read config file - error: %s\n", err.Error())
		utils.Print("RED", errMsg)
	} else {
		viper.Unmarshal(app.config)
	}

	app.expandHome()

	viper.OnConfigChange(func(in fsnotify.Event) {
		viper.Unmarshal(app.config)
		app.expandHome()
		if app.OnConfigChange != nil {
			app.OnConfigChange()
		}
	})
	viper.WatchConfig()
}

func (app *AppConfig) checkConfigErrors() {
	if app.configErr != nil {
		utils.Print("RED", "Couldn't get $XDG_CONFIG!")
		panic(app.configErr)
	}

	if app.cacheErr != nil {
		utils.Print("RED", "Couldn't get $XDG_CONFIG!")
		panic(app.cacheErr)
	}
}

func (app *AppConfig) expandHome() {
	app.config.CacheDir = utils.ExpandHomeDir(app.config.CacheDir)
}
