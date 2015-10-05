package main

import (
	"bytes"
	"strconv"
)

type TwitchApi struct {
	auth *TwitchAuth
}

func (api *TwitchApi) getChannel(channel string) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(channel)

	return ""
}

func (api *TwitchApi) getChannelVideos(channel string, limit int, offset int) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(channel)
	url.WriteString("/videos?limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return ""
}

func (api *TwitchApi) getChannelFollows(channel string, limit int, offset int) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(channel)
	url.WriteString("/follows?limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return ""
}

func (api *TwitchApi) getChannelTeams(channel string) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(channel)
	url.WriteString("/teams")

	return ""
}

func (api *TwitchApi) getChannelBadges(channel string) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/chat/")
	url.WriteString(channel)
	url.WriteString("/badges")

	return ""
}

func (api *TwitchApi) getEmotes() string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/chat/emoticons")

	return ""
}

func (api *TwitchApi) getUserFollows(user string, limit int, offset int) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/users/")
	url.WriteString(user)
	url.WriteString("/follows/channels?limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return ""
}

func (api *TwitchApi) isUserFollowing(user string, target string) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/users/")
	url.WriteString(user)
	url.WriteString("/follows/channels/")
	url.WriteString(target)

	return ""
}

func (api *TwitchApi) getGames(limit int, offset int) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/games/top?limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return ""
}

func (api *TwitchApi) searchChannels(query string, limit int, offset int) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/search/channels?q=")
	url.WriteString(query)
	url.WriteString("&limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return ""
}

func (api *TwitchApi) searchStreams(query string, limit int, offset int) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/search/streams?q=")
	url.WriteString(query)
	url.WriteString("&limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return ""
}

func (api *TwitchApi) searchGames(query string, queryType string, live bool) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/search/games?q=")
	url.WriteString(query)
	url.WriteString("&type=")
	url.WriteString(queryType)
	url.WriteString("&offset=")
	url.WriteString(strconv.FormatBool(live))

	return ""
}

func (api *TwitchApi) getStream(channel string) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/streams/")
	url.WriteString(channel)

	return ""
}

func (api *TwitchApi) getFeaturedStreams(limit int, offset int) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/streams/featured?limit=")
	url.WriteString(strconv.Itoa(limit))
	url.WriteString("&offset=")
	url.WriteString(strconv.Itoa(offset))

	return ""
}
