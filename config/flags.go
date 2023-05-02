package config

import "flag"

func (app *AppConfig) parseFlags() {
	flag.StringVar(&app.userConfigPath, "c", app.userConfigPath,
		"Specify The Directory to check for config.yml file.")
	flag.BoolVar(&app.config.HideImage, "hide-image", app.config.HideImage,
		"Do not display the cover art image.")
	flag.BoolVar(&app.config.RoundedCorners, "rounded-corners", app.config.RoundedCorners,
		"Enable Rounded Corners")
	flag.Parse()
}
