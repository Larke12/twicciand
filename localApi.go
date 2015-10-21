package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type LocalApi struct {
	YDLPath string
}

type ParamsUrlConv struct {
	Url string `json:"url"`
}

// Create a constructor so a new API object cannot be created without an auth key
func NewLocalApi(ydlPath string) *LocalApi {
	api := new(LocalApi)
	api.YDLPath = ydlPath

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
	result.Write(output)
	return result
}
