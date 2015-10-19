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
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type ParamsQueryType struct {
	Query     string `json:"query"`
	QueryType string `json:"query_type"`
	Live      bool   `json:"live"`
}

type ParamsTarget struct {
	Query  string `json:"query"`
	Target string `json:"target"`
}

type ParamsQuery struct {
	Query string `json:"query"`
}

type ParamsPage struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ParamsQueryFull struct {
	Query string     `json:"query"`
	Page  ParamsPage `json:"page_params"`
}

// This is the interface which describes the twitch API
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
func getApiUrl(url bytes.Buffer, api *twitchApi) bytes.Buffer {
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
	var data bytes.Buffer
	_, err = data.ReadFrom(response.Body)

	// Check if we read it correctly
	if err != nil {
		log.Print("Error receiving response from url:", url.String())
	}

	return data
}

// Returns a channel object, takes a ParamsQuery
func (api *twitchApi) getChannel(apiParams []byte) bytes.Buffer {
	params := new(ParamsQuery)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call getChannel")
	}

	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(params.Query)

	return getApiUrl(url, api)
}

func (api *twitchApi) getChannelVideos(apiParams []byte) bytes.Buffer {
	params := new(ParamsQueryFull)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call getChannelVideos")
	}

	var url bytes.Buffer

	// Compose the url for the request
	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(params.Query)
	url.WriteString("/videos?limit=")
	url.WriteString(strconv.Itoa(params.Page.Limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(params.Page.Offset))

	return getApiUrl(url, api)
}

func (api *twitchApi) getChannelFollows(apiParams []byte) bytes.Buffer {
	params := new(ParamsQueryFull)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call getChannelFollows")
	}

	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(params.Query)
	url.WriteString("/follows?limit=")
	url.WriteString(strconv.Itoa(params.Page.Limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(params.Page.Offset))

	return getApiUrl(url, api)
}

func (api *twitchApi) getChannelTeams(apiParams []byte) bytes.Buffer {
	params := new(ParamsQuery)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call getChannelTeams")
	}

	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(params.Query)
	url.WriteString("/teams")

	return getApiUrl(url, api)
}

func (api *twitchApi) getChannelBadges(apiParams []byte) bytes.Buffer {
	params := new(ParamsQuery)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call getChannelBadges")
	}

	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/chat/")
	url.WriteString(params.Query)
	url.WriteString("/badges")

	return getApiUrl(url, api)
}

func (api *twitchApi) getEmotes(apiParams []byte) bytes.Buffer {
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/chat/emoticons")

	return getApiUrl(url, api)
}

func (api *twitchApi) getUserFollows(apiParams []byte) bytes.Buffer {
	params := new(ParamsQueryFull)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call getUserFollows")
	}

	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/users/")
	url.WriteString(params.Query)
	url.WriteString("/follows/channels?limit=")
	url.WriteString(strconv.Itoa(params.Page.Limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(params.Page.Offset))

	return getApiUrl(url, api)
}

func (api *twitchApi) isUserFollowing(apiParams []byte) bytes.Buffer {
	params := new(ParamsTarget)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call isUserFollowing")
	}

	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/users/")
	url.WriteString(params.Query)
	url.WriteString("/follows/channels/")
	url.WriteString(params.Target)

	return getApiUrl(url, api)
}

func (api *twitchApi) getGames(apiParams []byte) bytes.Buffer {
	params := new(ParamsPage)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call getGames")
	}

	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/games/top?limit=")
	url.WriteString(strconv.Itoa(params.Limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(params.Offset))

	return getApiUrl(url, api)
}

func (api *twitchApi) searchChannels(apiParams []byte) bytes.Buffer {
	params := new(ParamsQueryFull)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call searchChannels")
	}

	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/search/channels?q=")
	url.WriteString(params.Query)
	url.WriteString("&limit=")
	url.WriteString(strconv.Itoa(params.Page.Limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(params.Page.Offset))

	return getApiUrl(url, api)
}

func (api *twitchApi) searchStreams(apiParams []byte) bytes.Buffer {
	params := new(ParamsQueryFull)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call searchStreams")
	}

	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/search/streams?q=")
	url.WriteString(params.Query)
	url.WriteString("&limit=")
	url.WriteString(strconv.Itoa(params.Page.Limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(params.Page.Offset))

	return getApiUrl(url, api)
}

func (api *twitchApi) searchGames(apiParams []byte) bytes.Buffer {
	params := new(ParamsQueryType)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call searchGames")
	}

	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/search/games?q=")
	url.WriteString(params.Query)
	url.WriteString("&type=")
	url.WriteString(params.QueryType)
	url.WriteString("&offset=")
	url.WriteString(strconv.FormatBool(params.Live))

	return getApiUrl(url, api)
}

func (api *twitchApi) getStream(apiParams []byte) bytes.Buffer {
	params := new(ParamsQuery)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call getStream")
	}

	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/streams/")
	url.WriteString(params.Query)

	return getApiUrl(url, api)
}

func (api *twitchApi) getFeaturedStreams(apiParams []byte) bytes.Buffer {
	params := new(ParamsPage)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Print("Incorrect parameters passed to twitch call getFeaturedStreams")
	}

	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/streams/featured?limit=")
	url.WriteString(strconv.Itoa(params.Limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(params.Offset))

	return getApiUrl(url, api)
}
