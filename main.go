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
	"path"
	"sync"


	"github.com/walle/cfg"
)

func main() {
	// Try finding the config file in the user's .config
	conffile := path.Join(os.Getenv("HOME"), ".config/twicciand/twicciand.conf")

	// If the config file doesn't exist, load from the same directory as the binary
	// (useful for testing)
	if _, err := os.Stat(conffile); err != nil {
		conffile = "twicciand.conf"
	}
	file, err := cfg.NewConfigFile(conffile)
	if err != nil {
		log.Print("Error parsing config file")
	}

	// Make a new authentication object
	auth := new(TwitchAuth)

	// Knowing the username is not necessary, but if it is provided, store it
	username, err := file.Config.GetString("username")
	if err != nil {
		log.Print("Could not read username")
		file.Config.SetString("username", "")
		file.Persist()
	}
	auth.Username = username

	// read the auth token from the config file, or receive it from twitch
	token, err := file.Config.GetString("token")
	if err != nil {
		log.Print("Could not find auth token - waiting for twitch's reply")
		// Wait until we receive the credentials
		auth.startAuthServer()
		// Update the config file
		file.Config.SetString("token", auth.Password)
		file.Persist()
	} else {
		// We have the pasword in the config file, inject it into the auth object
		auth.Password = token
	}

	// Print user's authentication token
	fmt.Println("Your token is:", auth.Password)

	twitchApi := NewTwitchApi(auth)
	// result := twitchApi.getChannelBadges([]byte(`{"query":"gamesdonequick"}`))
	// fmt.Println(result.String())

	// api := NewLocalApi("")
	// result = api.getStreamUrl([]byte(`{"url":"http://twitch.tv/stabbystabby"}`))
	// fmt.Println(result.String())

	reader := NewSocketReader(twitchApi)
	fmt.Println("Starting SocketReader...")
	var wg sync.WaitGroup
	wg.Add(1)
	go reader.StartReader()

	wg.Wait()
}
