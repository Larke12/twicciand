package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/gcfg.v1"
)

type Config struct {
	Twitch struct {
		Username string
		Token    string
	}
}

type TwitchAuth struct {
	Token  string
	Scopes string
}

// Handle our captive portal's post containing the token and allowed scopes
func handle_twitch_auth(com chan string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		com <- r.PostFormValue("token")
	}
}

func main() {
	cfg := new(Config)
	err := gcfg.ReadFileInto(cfg, "twicciand.conf")
	if err != nil {
		log.Fatalf("Failed to parse config data: %s", err)
	}

	if cfg.Twitch.Token == "" {
		com := make(chan string, 0)

		// Catch Twitch's authentication redirect which contains the token and list of scopes
		http.Handle("/", http.FileServer(http.Dir("auth_server")))
		// Recieve a post containing the token and list of scopes from our original capture page
		http.HandleFunc("/recv_auth", handle_twitch_auth(com))
		go http.ListenAndServe(":1921", nil)
		fmt.Println("Waiting for authentication token...")
		fmt.Println("Please visit", "https://api.twitch.tv/kraken/oauth2/authorize?response_type=token&client_id=mya9g4l7ucpsbwe2sjlj749d4hqzvvj&redirect_uri=http://localhost:1921&scope=user_read+user_follows_edit+user_subscriptions+chat_login", "to generate an authentication token")

		// Receive auth token from the channel
		cfg.Twitch.Token = <-com
	}

	// Print user's authentication token
	fmt.Println("Your token is:", cfg.Twitch.Token)
}
