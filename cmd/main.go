package main

import (
	"fmt"
	"log"

	"github.com/phuonganhniie/stp-terminal/auth"
	"github.com/phuonganhniie/stp-terminal/config"
)

func main() {
	// create and read configuration
	appConfig := config.NewAppConfig()
	appConfig.ReadConfig()

	// Authenticate the client
	authConfig := auth.NewAuthConfig(appConfig.UserConfigPath)
	spotifyClient, err := authConfig.InitClient()
	if err != nil {
		log.Fatal("Error initializing Spotify client: ", err)
	}

	// Your application logic with the authenticated Spotify client
	fmt.Printf("Spotify client authenticated: %+ v\n", spotifyClient)
}
