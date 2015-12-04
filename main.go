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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"

	"github.com/walle/cfg"
)

func main() {
	pid := os.Getpid()
	commandstr := "pgrep twicciand | grep -v " + strconv.Itoa(pid)
	out, _ := exec.Command("sh", "-c", commandstr).Output()
	if string(out) != "" {
		fmt.Println("The daemon is already running, terminating...")
		return
	}

	// Make a new authentication object
	auth := new(TwitchAuth)
	// Create new api objects
	twitchApi := NewTwitchApi(auth)

	// Run the socket reader
	reader := NewSocketReader(twitchApi)
	fmt.Println("Starting SocketReader...")
	var wg sync.WaitGroup
	wg.Add(1)
	go reader.StartReader()

	// Try finding the config file in the user's .config
	conffile := path.Join(os.Getenv("HOME"), ".config/twicciand/twicciand.conf")

	// If the config file doesn't exist, create one
	if _, err := os.Stat(path.Join(os.Getenv("HOME"), ".config/twicciand/")); os.IsNotExist(err) {
		os.Mkdir(path.Join(os.Getenv("HOME"), ".config/twicciand/"), 0755)
	}
	if _, err := os.Stat(conffile); err != nil {
		os.Create(conffile)
	}
	file, err := cfg.NewConfigFile(conffile)
	if err != nil {
		log.Print("Error parsing config file")
	}

	// read the auth token from the config file, or receive it from twitch
	token, err := file.Config.GetString("token")
	if err != nil || token == "" {
		log.Print("Could not find auth token - waiting for twitch's reply")
		// Wait until we receive the credentials
		auth.startAuthServer()
		// Update the config file
		file.Config.SetString("token", auth.Password)
	} else {
		// We have the pasword in the config file, inject it into the auth object
		auth.Password = token
	}

	// Knowing the username is not necessary, but if it is provided, store it
	username, err := file.Config.GetString("username")
	if err != nil || username == "" {
		result := twitchApi.getUser([]byte(`{"query":"nil"}`))
		var resultjson map[string]interface{}
		json.Unmarshal(result.Bytes(), &resultjson)

		auth.Username = resultjson["name"].(string)
		fmt.Println("Gotten username:", resultjson["name"].(string))

		file.Config.SetString("username", resultjson["name"].(string))
	}

	file.Persist()
	// Print user's authentication token
	fmt.Println("Your token is:", auth.Password)

	// Create a chat object
	chat := new(TwitchChat)
	chat.auth = auth
	chat.colorMap = make(map[string]string)
	chat.AddChannel(auth.Username, "#twitchplayspokemon", auth.Password)

	// Start chat server
	http.Handle("/ws", wsHandler{chat: chat})
	if err := http.ListenAndServe(":1922", nil); err != nil {
		log.Print("Error starting chat websocket server:", err)
	}

	wg.Wait()
}
