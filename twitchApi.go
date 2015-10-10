package main

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
)

type twitchApi struct {
	auth *TwitchAuth
}

// Create a constructor so a new API object cannot be created without an auth key
func NewTwitchApi(auth *TwitchAuth) *twitchApi {
	api := new(twitchApi)
	api.auth = auth

	return api
}

// Take a URL and make a GET request to twitch's REST api
func getJsonRequest(url bytes.Buffer, api *twitchApi) bytes.Buffer {
	// Create a HTTP request
	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Set("Accept", "application/vnd.twitchtv.v3+json") // Request the v3 api
	req.Header.Set("Client-ID", api.auth.Password)

	// Run that request
	client := new(http.Client)
	response, err := client.Do(req)
	if err != nil {
		log.Print("Error making GET request to url:", url.String())
	}

	// Capture output in a bytes.Buffer
	var json bytes.Buffer
	_, err = json.ReadFrom(response.Body)

	// Check if we read it correctly
	if err != nil {
		log.Print("Error receiving response from url:", url.String())
	}

	return json
}

func (api *twitchApi) getChannel(channel string) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(channel)

	return getJsonRequest(url, api)
}

func (api *twitchApi) getChannelVideos(channel string, limit int, offset int) bytes.Buffer {
	var url bytes.Buffer

	// Compose the url for the request
	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(channel)
	url.WriteString("/videos?limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return getJsonRequest(url, api)
}

func (api *twitchApi) getChannelFollows(channel string, limit int, offset int) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(channel)
	url.WriteString("/follows?limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return getJsonRequest(url, api)
}

func (api *twitchApi) getChannelTeams(channel string) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(channel)
	url.WriteString("/teams")

	return getJsonRequest(url, api)
}

func (api *twitchApi) getChannelBadges(channel string) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/chat/")
	url.WriteString(channel)
	url.WriteString("/badges")

	return getJsonRequest(url, api)
}

func (api *twitchApi) getEmotes() bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/chat/emoticons")

	return getJsonRequest(url, api)
}

func (api *twitchApi) getUserFollows(user string, limit int, offset int) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/users/")
	url.WriteString(user)
	url.WriteString("/follows/channels?limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return getJsonRequest(url, api)
}

func (api *twitchApi) isUserFollowing(user string, target string) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/users/")
	url.WriteString(user)
	url.WriteString("/follows/channels/")
	url.WriteString(target)

	return getJsonRequest(url, api)
}

func (api *twitchApi) getGames(limit int, offset int) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/games/top?limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return getJsonRequest(url, api)
}

func (api *twitchApi) searchChannels(query string, limit int, offset int) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/search/channels?q=")
	url.WriteString(query)
	url.WriteString("&limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return getJsonRequest(url, api)
}

func (api *twitchApi) searchStreams(query string, limit int, offset int) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/search/streams?q=")
	url.WriteString(query)
	url.WriteString("&limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return getJsonRequest(url, api)
}

func (api *twitchApi) searchGames(query string, queryType string, live bool) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/search/games?q=")
	url.WriteString(query)
	url.WriteString("&type=")
	url.WriteString(queryType)
	url.WriteString("&offset=")
	url.WriteString(strconv.FormatBool(live))

	return getJsonRequest(url, api)
}

func (api *twitchApi) getStream(channel string) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/streams/")
	url.WriteString(channel)

	return getJsonRequest(url, api)
}

func (api *twitchApi) getFeaturedStreams(limit int, offset int) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/streams/featured?limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return getJsonRequest(url, api)
}
