package main

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/gcfg.v1"
)

func TestGetChannel(t *testing.T) {
	cfg := new(Config)
	err := gcfg.ReadFileInto(cfg, "twicciand.conf")
	if err != nil {
		log.Printf("Failed to parse config data: %s", err)
	}

	auth := new(TwitchAuth)

	// Create new authentication storage
	if cfg.Twitch.Token == "" {
		// Whelp, return since we can't run any tests
		return
	} else {
		// We have the pasword in the config file
		auth.setPassword(cfg.Twitch.Token)
	}

	api := NewTwitchApi(auth)

	// Start testing
	Convey("Test API results for getChannelVideos", t, func() {
		var expected bytes.Buffer

		Convey("Test twitch's test_channel", func() {
			expected.WriteString(`{"_total":9,"_links":{"self":"https://api.twitch.tv/kraken/channels/test_channel/videos?limit=1&offset=0&user=test_channel","next":"https://api.twitch.tv/kraken/channels/test_channel/videos?limit=1&offset=1&user=test_channel"},"videos":[{"title":"robot greeting 2","description":"greeting","broadcast_id":null,"status":"recorded","tag_list":"","_id":"c213462","recorded_at":"2009-12-15T10:03:04Z","game":null,"length":75,"is_muted":false,"preview":null,"url":"http://www.twitch.tv/test_channel/c/213462","views":2,"fps":null,"resolutions":null,"broadcast_type":"highlight","created_at":"2009-12-15T10:03:04Z","_links":{"self":"https://api.twitch.tv/kraken/videos/c213462","channel":"https://api.twitch.tv/kraken/channels/test_channel"},"channel":{"name":"test_channel","display_name":"Test_channel"}}]}`)
			result := api.getChannelVideos("test_channel", 1, 0)

			So(expected, ShouldResemble, result)
		})
		Convey("Test gamesdonequick's channel", func() {
			result := api.getChannelVideos("gamesdonequick", 10, 0)

			resultjson := TwitchChannelVideos{}
			json.Unmarshal(result.Bytes(), &resultjson)
			So(len(resultjson.Videos), ShouldBeGreaterThan, 0)
		})
	})

}
