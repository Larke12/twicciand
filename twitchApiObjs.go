package main

type TwitchChannel struct {
	Id    int `json:"_id"`
	Links struct {
		Chat          string `json:"chat"`
		Commercial    string `json:"commercial"`
		Editors       string `json:"editors"`
		Features      string `json:"features"`
		Follows       string `json:"follows"`
		Self          string `json:"self"`
		StreamKey     string `json:"stream_key"`
		Subscriptions string `json:"subscriptions"`
		Teams         string `json:"teams"`
		Videos        string `json:"videos"`
	} `json:"_links"`
	Background                   interface{} `json:"background"`
	Banner                       interface{} `json:"banner"`
	BroadcasterLanguage          string      `json:"broadcaster_language"`
	CreatedAt                    string      `json:"created_at"`
	Delay                        int         `json:"delay"`
	DisplayName                  string      `json:"display_name"`
	Followers                    int         `json:"followers"`
	Game                         string      `json:"game"`
	Language                     string      `json:"language"`
	Logo                         string      `json:"logo"`
	Mature                       bool        `json:"mature"`
	Name                         string      `json:"name"`
	Partner                      bool        `json:"partner"`
	ProfileBanner                interface{} `json:"profile_banner"`
	ProfileBannerBackgroundColor interface{} `json:"profile_banner_background_color"`
	Status                       string      `json:"status"`
	UpdatedAt                    string      `json:"updated_at"`
	URL                          string      `json:"url"`
	VideoBanner                  string      `json:"video_banner"`
	Views                        int         `json:"views"`
}

type TwitchChannelVideos struct {
	Links struct {
		Next string `json:"next"`
		Self string `json:"self"`
	} `json:"_links"`
	Total  int `json:"_total"`
	Videos []struct {
		Id    string `json:"_id"`
		Links struct {
			Channel string `json:"channel"`
			Self    string `json:"self"`
		} `json:"_links"`
		BroadcastID   interface{} `json:"broadcast_id"`
		BroadcastType string      `json:"broadcast_type"`
		Channel       struct {
			DisplayName string `json:"display_name"`
			Name        string `json:"name"`
		} `json:"channel"`
		CreatedAt   string      `json:"created_at"`
		Description string      `json:"description"`
		Fps         interface{} `json:"fps"`
		Game        interface{} `json:"game"`
		IsMuted     bool        `json:"is_muted"`
		Length      int         `json:"length"`
		Preview     interface{} `json:"preview"`
		RecordedAt  string      `json:"recorded_at"`
		Resolutions interface{} `json:"resolutions"`
		Status      string      `json:"status"`
		TagList     string      `json:"tag_list"`
		Title       string      `json:"title"`
		URL         string      `json:"url"`
		Views       int         `json:"views"`
	} `json:"videos"`
}

type TwitchChannelFollows struct {
	Links struct {
		Next string `json:"next"`
		Self string `json:"self"`
	} `json:"_links"`
	Total   int `json:"_total"`
	Follows []struct {
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
		CreatedAt     string `json:"created_at"`
		Notifications bool   `json:"notifications"`
		User          struct {
			Id    int `json:"_id"`
			Links struct {
				Self string `json:"self"`
			} `json:"_links"`
			Bio         interface{} `json:"bio"`
			CreatedAt   string      `json:"created_at"`
			DisplayName string      `json:"display_name"`
			Logo        string      `json:"logo"`
			Name        string      `json:"name"`
			Type        string      `json:"type"`
			UpdatedAt   string      `json:"updated_at"`
		} `json:"user"`
	} `json:"follows"`
}

type TwitchChannelTeams struct {
	Links struct {
		Self string `json:"self"`
	} `json:"_links"`
	Teams []struct {
		Id    int `json:"_id"`
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
		Background  interface{} `json:"background"`
		Banner      string      `json:"banner"`
		CreatedAt   string      `json:"created_at"`
		DisplayName string      `json:"display_name"`
		Info        string      `json:"info"`
		Logo        string      `json:"logo"`
		Name        string      `json:"name"`
		UpdatedAt   string      `json:"updated_at"`
	} `json:"teams"`
}

type TwitchChannelBadges struct {
	Links struct {
		Self string `json:"self"`
	} `json:"_links"`
	Admin struct {
		Alpha string `json:"alpha"`
		Image string `json:"image"`
		Svg   string `json:"svg"`
	} `json:"admin"`
	Broadcaster struct {
		Alpha string `json:"alpha"`
		Image string `json:"image"`
		Svg   string `json:"svg"`
	} `json:"broadcaster"`
	GlobalMod struct {
		Alpha string `json:"alpha"`
		Image string `json:"image"`
		Svg   string `json:"svg"`
	} `json:"global_mod"`
	Mod struct {
		Alpha string `json:"alpha"`
		Image string `json:"image"`
		Svg   string `json:"svg"`
	} `json:"mod"`
	Staff struct {
		Alpha string `json:"alpha"`
		Image string `json:"image"`
		Svg   string `json:"svg"`
	} `json:"staff"`
	Subscriber struct {
		Image string `json:"image"`
	} `json:"subscriber"`
	Turbo struct {
		Alpha string `json:"alpha"`
		Image string `json:"image"`
		Svg   string `json:"svg"`
	} `json:"turbo"`
}
