package config

import "flag"

func (app *AppConfig) parseFlags() {
	flag.StringVar(&app.UserConfigPath, "c", app.UserConfigPath,
		"Specify The Directory to check for config.yaml file.")
	flag.BoolVar(&app.config.HideImage, "hide-image", app.config.HideImage,
		"Do not display the cover art image.")
	flag.BoolVar(&app.config.RoundedCorners, "rounded-corners", app.config.RoundedCorners,
		"Enable Rounded Corners")
	flag.Parse()
}
