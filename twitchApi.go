package main

import (
	"bytes"
)

type TwitchApi struct {
	auth *TwitchAuth
}

type ChannelTemplate struct {
	Channel string
}

func (api *TwitchApi) getChannel(channel string) string {
	// TODO: implement
	var url bytes.Buffer

	url.WriteString("https://api.twitch.tv/kraken/channels/")
	url.WriteString(channel)

	return ""
}

func (api *TwitchApi) getChannelVideos(channel string) string {
	// TODO: implement
	return ""
}

func (api *TwitchApi) getChannelFollows(channel string) string {
	// TODO: implement
	return ""
}

func (api *TwitchApi) getChannelTeams(channel string) string {
	// TODO: implement
	return ""
}

func (api *TwitchApi) getChannelBadges(channel string) string {
	// TODO: implement
	return ""
}

func (api *TwitchApi) getEmotes() string {
	// TODO: implement
	return ""
}

func (api *TwitchApi) getUserFollows(user string) string {
	// TODO: implement
	return ""
}

func (api *TwitchApi) isUserFollowing(user string, target string) string {
	// TODO: implement
	return ""
}

func (api *TwitchApi) getGames() string {
	// TODO: implement
	return ""
}

func (api *TwitchApi) searchChannels(query string) string {
	// TODO: implement
	return ""
}

func (api *TwitchApi) searchStreams(query string) string {
	// TODO: implement
	return ""
}

func (api *TwitchApi) searchGames(query string) string {
	// TODO: implement
	return ""
}

func (api *TwitchApi) getStream(query string) string {
	// TODO: implement
	return ""
}

func (api *TwitchApi) getFeaturedStreams() string {
	// TODO: implement
	return ""
}
