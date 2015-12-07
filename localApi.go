package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type LocalApi struct {
	YDLPath string
	auth    *TwitchAuth
	chat    *TwitchChat
}

type ParamsUrlConv struct {
	Url string `json:"url"`
}

type ParamsLocal struct {
	Query string `json:"query"`
}

// Create a constructor so a new API object cannot be created without an auth key
func NewLocalApi(ydlPath string, auth *TwitchAuth, chat *TwitchChat) *LocalApi {
	api := new(LocalApi)
	api.YDLPath = ydlPath
	api.auth = auth
	api.chat = chat

	return api
}

// Gets the actual stream URL using youtube-dl
func (api *LocalApi) getStreamUrl(apiParams []byte) bytes.Buffer {
	fmt.Println(string(apiParams))
	params := new(ParamsUrlConv)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Printf("Incorrect parameters passed to local call getStreamUrl: %s", err)
	}

	// Capture youtube-dl output
	var output []byte
	args := []string{"-g", params.Url}
	if output, err = exec.Command("youtube-dl", args...).Output(); err != nil {
		log.Printf("There was a problem running youtube-dl: %s", err)
	}

	var result bytes.Buffer
	result.WriteString(`"`)
	result.WriteString(strings.TrimSpace(string(output)))
	result.WriteString(`"`)
	return result
}

// Gets the actual stream URL using youtube-dl
func (api *LocalApi) getStreamDesc(apiParams []byte) bytes.Buffer {
	// fmt.Println(string(apiParams))
	params := new(ParamsUrlConv)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Printf("Incorrect parameters passed to local call getStreamUrl: %s", err)
	}

	// Capture youtube-dl output
	var output []byte
	args := []string{"--get-description", params.Url}
	if output, err = exec.Command("youtube-dl", args...).Output(); err != nil {
		log.Printf("There was a problem running youtube-dl: %s", err)
	}

	var result bytes.Buffer
	result.WriteString(`"`)
	result.Write(output)
	result.WriteString(`"`)
	return result
}

// Gets the actual stream URL using youtube-dl
func (api *LocalApi) isAuthenticated(apiParams []byte) bytes.Buffer {
	var result bytes.Buffer
	if api.auth.Password == "" {
		result.WriteString("false")
	} else {
		result.WriteString("true")
	}
	return result
}

// Changes the current chat channel
func (api *LocalApi) changeChat(apiParams []byte) bytes.Buffer {
	params := new(ParamsLocal)
	err := json.Unmarshal(apiParams, params)
	if err != nil {
		log.Printf("There was a problem parsing the requested chat channel")
	}

	api.chat.AddChannel(api.auth.Username, "#"+params.Query, api.auth.Password)

	var result bytes.Buffer
	result.WriteString("true")
	return result
}
