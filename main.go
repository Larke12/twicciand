package main

import (
	"fmt"
	"log"

	"gopkg.in/gcfg.v1"
)

type Config struct {
	Twitch struct {
		Username string
		Token    string
	}
}


func main() {
	cfg := new(Config)
	err := gcfg.ReadFileInto(cfg, "twicciand.conf")
	if err != nil {
		log.Fatalf("Failed to parse config data: %s", err)
	}

	auth := new(TwitchAuth)

	// Knowing the username is not necessary, but if it is provided, store it
	if cfg.Twitch.Username != "" {
		auth.setUsername(cfg.Twitch.Username)
	}

	// Create new authentication storage
	if cfg.Twitch.Token == "" {
		// Wait until we receive the credentials
		auth.startAuthServer()
	} else {
		// We have the pasword in the config file
		auth.setPassword(cfg.Twitch.Token)
	}

	// Print user's authentication token
	fmt.Println("Your token is:", auth.Password)
}
