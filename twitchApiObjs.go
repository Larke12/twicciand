package main

type TwitchChannelVideos struct {
	Links struct {
		Next string `json:"next"`
		Self string `json:"self"`
	} `json:"_links"`
	Total  int `json:"_total"`
	Videos []struct {
		_id    string `json:"_id"`
		_links struct {
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
