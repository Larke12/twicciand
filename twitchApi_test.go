package main

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/walle/cfg"
)

func TestChannel(t *testing.T) {
	file, err := cfg.NewConfigFile("twicciand.conf")
	if err != nil {
		log.Printf("Failed to parse config data: %s", err)
	}

	auth := new(TwitchAuth)

	username, err := file.Config.GetString("username")
	if err != nil {
		log.Print("Could not read username")
		file.Config.SetString("username", "")
		file.Persist()
	}
	auth.Username = username

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

	api := NewTwitchApi(auth)

	// Start testing
	Convey("Test API results for getChannelVideos", t, func() {
		var expected bytes.Buffer

		Convey("Test twitch's test_channel", func() {
			expected.WriteString(`{"_total":9,"_links":{"self":"https://api.twitch.tv/kraken/channels/test_channel/videos?limit=1&offset=0&user=test_channel","next":"https://api.twitch.tv/kraken/channels/test_channel/videos?limit=1&offset=1&user=test_channel"},"videos":[{"title":"robot greeting 2","description":"greeting","broadcast_id":null,"status":"recorded","tag_list":"","_id":"c213462","recorded_at":"2009-12-15T10:03:04Z","game":null,"length":75,"is_muted":false,"preview":null,"url":"http://www.twitch.tv/test_channel/c/213462","views":2,"fps":null,"resolutions":null,"broadcast_type":"highlight","created_at":"2009-12-15T10:23:51Z","_links":{"self":"https://api.twitch.tv/kraken/videos/c213462","channel":"https://api.twitch.tv/kraken/channels/test_channel"},"channel":{"name":"test_channel","display_name":"Test_channel"}}]}`)
			result := api.getChannelVideos([]byte(`{"query":"test_channel","page_params":{"limit":1,"offset":0}}`))

			So(expected.String(), ShouldResemble, result.String())
		})
		Convey("Test a property of a real channel", func() {
			result := api.getChannelVideos([]byte(`{"query":"gamesdonequick","page_params":{"limit":10,"offset":0}}`))

			resultjson := TwitchChannelVideos{}
			json.Unmarshal(result.Bytes(), &resultjson)
			So(len(resultjson.Videos), ShouldBeGreaterThan, 0)
		})
	})

	Convey("Test API results for getChannel", t, func() {
		Convey("Test twitch's test_channel", func() {
			result := api.getChannel([]byte(`{"query":"test_channel"}`))

			var resultjson map[string]interface{}
			json.Unmarshal(result.Bytes(), &resultjson)
			So(len(resultjson["name"].(string)), ShouldBeGreaterThan, 0)
		})
	})

	Convey("Test API results for getChannelFollows", t, func() {
		Convey("Test twitch's test_channel", func() {
			result := api.getChannelFollows([]byte(`{"query":"test_channel","page_params":{"limit":1,"offset":0}}`))

			var resultjson map[string]interface{}
			json.Unmarshal(result.Bytes(), &resultjson)
			So(len(resultjson["follows"].([]interface{})), ShouldBeGreaterThan, 0)
		})
	})

	Convey("Test API results for getChannelTeams", t, func() {
		var expected bytes.Buffer

		Convey("Test twitch's test_channel", func() {
			expected.WriteString(`{"_links":{"self":"https://api.twitch.tv/kraken/channels/test_channel/teams"},"teams":[]}`)
			result := api.getChannelTeams([]byte(`{"query":"test_channel","page_params":{"limit":1,"offset":0}}`))

			So(expected.String(), ShouldResemble, result.String())
		})
	})

}

func TestChat(t *testing.T) {
	file, err := cfg.NewConfigFile("twicciand.conf")
	if err != nil {
		log.Printf("Failed to parse config data: %s", err)
	}

	auth := new(TwitchAuth)

	username, err := file.Config.GetString("username")
	if err != nil {
		log.Print("Could not read username")
		file.Config.SetString("username", "")
		file.Persist()
	}
	auth.Username = username

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

	api := NewTwitchApi(auth)

	// Start testing
	Convey("Test API results for getEmotes", t, func() {
		Convey("Test twitch's emotes endpoint", func() {
			result := api.getEmotes([]byte(``))

			var resultjson map[string]interface{}
			json.Unmarshal(result.Bytes(), &resultjson)
			So(len(resultjson["emoticons"].([]interface{})), ShouldBeGreaterThan, 0)
		})
	})

	Convey("Test API results for getChannelBadges", t, func() {
		var expected bytes.Buffer

		Convey("Test twitch's test_channel", func() {
			expected.WriteString(`{"global_mod":{"alpha":"http://chat-badges.s3.amazonaws.com/globalmod-alpha.png","image":"http://chat-badges.s3.amazonaws.com/globalmod.png","svg":"http://chat-badges.s3.amazonaws.com/globalmod.svg"},"admin":{"alpha":"http://chat-badges.s3.amazonaws.com/admin-alpha.png","image":"http://chat-badges.s3.amazonaws.com/admin.png","svg":"http://chat-badges.s3.amazonaws.com/admin.svg"},"broadcaster":{"alpha":"http://chat-badges.s3.amazonaws.com/broadcaster-alpha.png","image":"http://chat-badges.s3.amazonaws.com/broadcaster.png","svg":"http://chat-badges.s3.amazonaws.com/broadcaster.svg"},"mod":{"alpha":"http://chat-badges.s3.amazonaws.com/mod-alpha.png","image":"http://chat-badges.s3.amazonaws.com/mod.png","svg":"http://chat-badges.s3.amazonaws.com/mod.svg"},"staff":{"alpha":"http://chat-badges.s3.amazonaws.com/staff-alpha.png","image":"http://chat-badges.s3.amazonaws.com/staff.png","svg":"http://chat-badges.s3.amazonaws.com/staff.svg"},"turbo":{"alpha":"http://chat-badges.s3.amazonaws.com/turbo-alpha.png","image":"http://chat-badges.s3.amazonaws.com/turbo.png","svg":"http://chat-badges.s3.amazonaws.com/turbo.svg"},"subscriber":null,"_links":{"self":"https://api.twitch.tv/kraken/chat/test_channel/badges"}}`)
			result := api.getChannelBadges([]byte(`{"query":"test_channel"}`))

			So(expected.String(), ShouldResemble, result.String())
		})
	})

}

func TestUser(t *testing.T) {
	file, err := cfg.NewConfigFile("twicciand.conf")
	if err != nil {
		log.Printf("Failed to parse config data: %s", err)
	}

	auth := new(TwitchAuth)

	username, err := file.Config.GetString("username")
	if err != nil {
		log.Print("Could not read username")
		file.Config.SetString("username", "")
		file.Persist()
	}
	auth.Username = username

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

	api := NewTwitchApi(auth)

	// Start testing
	Convey("Test API results for getUserFollows", t, func() {
		var expected bytes.Buffer

		Convey("Test twitch's test_user1", func() {
			expected.WriteString(`{"follows":[],"_total":0,"_links":{"self":"https://api.twitch.tv/kraken/users/test_user1/follows/channels?direction=DESC&limit=1&offset=0&sortby=created_at","next":"https://api.twitch.tv/kraken/users/test_user1/follows/channels?direction=DESC&limit=1&offset=1&sortby=created_at"}}`)
			result := api.getUserFollows([]byte(`{"query":"test_user1","page_params":{"limit":1,"offset":0}}`))

			So(expected.String(), ShouldResemble, result.String())
		})
	})

	Convey("Test API results for isUserFollowing", t, func() {
		var expected bytes.Buffer

		Convey("Test if finaleti is following crumps2", func() {
			expected.WriteString(`{"created_at":"2013-07-20T04:09:33+00:00","_links":{"self":"https://api.twitch.tv/kraken/users/finaleti/follows/channels/crumps2"},"notifications":true,"channel":{"mature":true,"status":"MLb 15 - Dirtbag can't catch a break","broadcaster_language":"en","display_name":"Crumps2","game":"MLB 15: The Show","delay":0,"language":"en","_id":19107317,"name":"crumps2","created_at":"2010-12-29T22:02:50Z","updated_at":"2015-10-16T00:16:38Z","logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/crumps2-profile_image-40d32b958f59a0c5-300x300.jpeg","banner":null,"video_banner":"http://static-cdn.jtvnw.net/jtv_user_pictures/crumps2-channel_offline_image-2fac52e223148bd1-1920x1080.jpeg","background":null,"profile_banner":"http://static-cdn.jtvnw.net/jtv_user_pictures/crumps2-profile_banner-2ccbe7d1eb2197fb-480.png","profile_banner_background_color":null,"partner":true,"url":"http://www.twitch.tv/crumps2","views":9485484,"followers":77442,"_links":{"self":"https://api.twitch.tv/kraken/channels/crumps2","follows":"https://api.twitch.tv/kraken/channels/crumps2/follows","commercial":"https://api.twitch.tv/kraken/channels/crumps2/commercial","stream_key":"https://api.twitch.tv/kraken/channels/crumps2/stream_key","chat":"https://api.twitch.tv/kraken/chat/crumps2","features":"https://api.twitch.tv/kraken/channels/crumps2/features","subscriptions":"https://api.twitch.tv/kraken/channels/crumps2/subscriptions","editors":"https://api.twitch.tv/kraken/channels/crumps2/editors","teams":"https://api.twitch.tv/kraken/channels/crumps2/teams","videos":"https://api.twitch.tv/kraken/channels/crumps2/videos"}}}`)
			result := api.isUserFollowing([]byte(`{"query":"finaleti","target":"crumps2"}`))

			So(expected.String(), ShouldResemble, result.String())
		})
	})

}

func TestGames(t *testing.T) {
	file, err := cfg.NewConfigFile("twicciand.conf")
	if err != nil {
		log.Printf("Failed to parse config data: %s", err)
	}

	auth := new(TwitchAuth)

	username, err := file.Config.GetString("username")
	if err != nil {
		log.Print("Could not read username")
		file.Config.SetString("username", "")
		file.Persist()
	}
	auth.Username = username

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

	api := NewTwitchApi(auth)

	// Start testing
	Convey("Test API results for getGames", t, func() {
		Convey("Test for top games on Twitch", func() {
			result := api.getGames([]byte(`{"page_params":{"limit":1,"offset":0}}`))

			var resultjson map[string]interface{}
			json.Unmarshal(result.Bytes(), &resultjson)
			So(len(resultjson["top"].([]interface{})), ShouldBeGreaterThan, 0)
		})
	})

	Convey("Test API results for searchChannels", t, func() {
		Convey("Test for starcraft channels", func() {
			result := api.searchChannels([]byte(`{"query":"starcraft","page_params":{"limit":1,"offset":0}}`))

			var resultjson map[string]interface{}
			json.Unmarshal(result.Bytes(), &resultjson)
			So(len(resultjson["channels"].([]interface{})), ShouldBeGreaterThan, 0)
		})
	})

	Convey("Test API results for searchStreams", t, func() {
		Convey("Test for starcraft streams", func() {
			result := api.searchStreams([]byte(`{"query":"starcraft","page_params":{"limit":1,"offset":0}}`))

			var resultjson map[string]interface{}
			json.Unmarshal(result.Bytes(), &resultjson)
			So(len(resultjson["streams"].([]interface{})), ShouldBeGreaterThan, 0)
		})
	})

	Convey("Test API results for searchGames", t, func() {
		Convey("Test for starcraft games", func() {
			result := api.searchGames([]byte(`{"query":"star","query_type":"suggest","live":true}`))

			var resultjson map[string]interface{}
			json.Unmarshal(result.Bytes(), &resultjson)
			So(len(resultjson["games"].([]interface{})), ShouldBeGreaterThan, 0)
		})
	})

	Convey("Test API results for getStream", t, func() {
		Convey("Test for Twitch's test channel", func() {
			result := api.getStream([]byte(`{"query":"twitchplayspokemon"}`))

			var resultjson map[string]interface{}
			json.Unmarshal(result.Bytes(), &resultjson)
			So(len(resultjson["_links"].(map[string]interface{})), ShouldBeGreaterThan, 0)
		})
	})

	Convey("Test API results for getFeaturedStreams", t, func() {
		Convey("Test for list of featured streams", func() {
			result := api.getFeaturedStreams([]byte(`{"page_params":{"limit":1,"offset":0}}`))

			var resultjson map[string]interface{}
			json.Unmarshal(result.Bytes(), &resultjson)
			So(len(resultjson["featured"].([]interface{})), ShouldBeGreaterThan, 0)
		})
	})

	Convey("Test API results for getFollowedStreams", t, func() {
		Convey("Test for list of followed streams", func() {
			result := api.getFollowedStreams([]byte(`{"page_params":{"limit":1,"offset":0}}`))

			var resultjson map[string]interface{}
			json.Unmarshal(result.Bytes(), &resultjson)
			So(len(resultjson["streams"].([]interface{})), ShouldBeGreaterThan, 0)
		})
	})
}
