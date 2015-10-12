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
	"net/http"
)

// Generic Authentication provider interface
type Auth interface {
	isAuthenticated() bool
	setUsername(user string)
	setPassword(pass string)
	setCredentials(user string, pass string)
	getUsername() string
	getPassword() string
}

// Create a type for Twitch authentication
type TwitchAuth struct {
	Username string
	Password string
}

// Check if a authentication object has credentials stored
func (auth *TwitchAuth) isAuthenticated() bool {
	if auth.Username != "" && auth.Password != "" {
		return true
	} else {
		return false
	}
}

// Store credentials into the authentication object
func (auth *TwitchAuth) setUsername(user string) {
	auth.Username = user
}

func (auth *TwitchAuth) setPassword(pass string) {
	auth.Password = pass
}

func (auth *TwitchAuth) setCredentials(user string, pass string) {
	auth.Username = user
	auth.Password = pass
}

// Retrieve credentials from the authentication object
func (auth *TwitchAuth) getUsername() string {
	return auth.Username
}

func (auth *TwitchAuth) getPassword() string {
	return auth.Password
}

// Below are the functions to create a webserver to recieve credentials from Twitch

// Handle our captive portal's post containing the token and allowed scopes
func handle_twitch_auth(com chan string) http.HandlerFunc {
	// Return a function which will receive the posted value and push it down the channel
	return func(w http.ResponseWriter, r *http.Request) {
		com <- r.PostFormValue("token")
	}
}

// Start the webserver and block until we get credentials
func (auth *TwitchAuth) startAuthServer() {
	// Create a channel to pass to the webserver handler
	com := make(chan string, 0)

	// Catch Twitch's authentication redirect which contains the token and list of scopes
	http.Handle("/", http.FileServer(http.Dir("auth_server")))
	// Recieve a post containing the token and list of scopes from our original capture page
	http.HandleFunc("/recv_auth", handle_twitch_auth(com))
	go http.ListenAndServe(":19210", nil)

	// Print instructions
	fmt.Println("Waiting for authentication token...")
	fmt.Println("Please visit", "https://api.twitch.tv/kraken/oauth2/authorize?response_type=token&client_id=mya9g4l7ucpsbwe2sjlj749d4hqzvvj&redirect_uri=http://localhost:19210&scope=user_read+user_follows_edit+user_subscriptions+chat_login", "to generate an authentication token")

	// Receive auth token from the channel
	auth.Password = <-com
}
