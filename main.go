// This file is part of Twicciand.
// 
// Twicciand is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// Twicciand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with Twicciand.  If not, see <http://www.gnu.org/licenses/>.

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

	api := NewTwitchApi(auth)
	result := api.getChannelVideos("gamesdonequick", 5, 0)
	fmt.Println(result.String())
}
