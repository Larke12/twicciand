package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/gcfg.v1"
)

type Config struct {
	Twitch struct {
		Username string
		Token    string
	}
}

func writeConfig(auth Auth) {
		file, err := os.OpenFile("twicciand.conf", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Printf("Could not open config file to add auth token")
		}

		defer file.Close()

		file.WriteString("[twitch]\nusername=")
		file.WriteString(auth.getUsername())
		file.WriteString("\ntoken=")
		file.WriteString(auth.getPassword())
		file.WriteString("\n")

		file.Sync()
}

func main() {
	cfg := new(Config)
	err := gcfg.ReadFileInto(cfg, "twicciand.conf")
	if err != nil {
		log.Printf("Failed to parse config data: %s", err)
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
		writeConfig(auth)
	} else {
		// We have the pasword in the config file
		auth.setPassword(cfg.Twitch.Token)
	}

	// Print user's authentication token
	fmt.Println("Your token is:", auth.Password)
}
