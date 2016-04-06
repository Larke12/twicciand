package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type JsonRpc struct {
	Api    string                 `json:"api"`
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params"`
}

type JsonRpcResult struct {
	Name   string      `json:"name"`
	Result interface{} `json:"result"`
}

type SocketReader struct {
	Twitch        *TwitchApi
	Local         *LocalApi
	Listener      net.Listener
	TwitchFuncmap map[string]func(*TwitchApi, []byte) bytes.Buffer
	LocalFuncmap  map[string]func(*LocalApi, []byte) bytes.Buffer
}

// Properly create a new socket reader
func NewSocketReader(api *TwitchApi, chat *TwitchChat) *SocketReader {
	read := new(SocketReader)
	read.Twitch = api
	read.Local = NewLocalApi("", api.auth, chat)

	ln, err := net.Listen("tcp", ":1921")
	if err != nil {
		log.Panic("Could not open socketReader")
	}
	read.Listener = ln

	// Load twitch api functions into a function map, so we can dispatch calls easily
	read.TwitchFuncmap = make(map[string]func(*TwitchApi, []byte) bytes.Buffer)
	read.TwitchFuncmap["getChannel"] = (*TwitchApi).getChannel
	read.TwitchFuncmap["getChannelVideos"] = (*TwitchApi).getChannelVideos
	read.TwitchFuncmap["getChannelFollows"] = (*TwitchApi).getChannelFollows
	read.TwitchFuncmap["getChannelTeams"] = (*TwitchApi).getChannelTeams
	read.TwitchFuncmap["getChannelBadges"] = (*TwitchApi).getChannelBadges
	read.TwitchFuncmap["getEmotes"] = (*TwitchApi).getEmotes
	read.TwitchFuncmap["getUserFollows"] = (*TwitchApi).getUserFollows
	read.TwitchFuncmap["isUserFollowing"] = (*TwitchApi).isUserFollowing
	read.TwitchFuncmap["getGames"] = (*TwitchApi).getGames
	read.TwitchFuncmap["searchChannels"] = (*TwitchApi).searchChannels
	read.TwitchFuncmap["searchStreams"] = (*TwitchApi).searchStreams
	read.TwitchFuncmap["searchGames"] = (*TwitchApi).searchGames
	read.TwitchFuncmap["getStream"] = (*TwitchApi).getStream
	read.TwitchFuncmap["getFeaturedStreams"] = (*TwitchApi).getFeaturedStreams
	read.TwitchFuncmap["getFollowedStreams"] = (*TwitchApi).getFollowedStreams

	// Load local api functions into a function map, so we can dispatch calls easily
	read.LocalFuncmap = make(map[string]func(*LocalApi, []byte) bytes.Buffer)
	read.LocalFuncmap["getStreamUrl"] = (*LocalApi).getStreamUrl
	read.LocalFuncmap["getStreamDesc"] = (*LocalApi).getStreamDesc
	read.LocalFuncmap["changeChat"] = (*LocalApi).changeChat
	read.LocalFuncmap["isAuthenticated"] = (*LocalApi).isAuthenticated

	return read
}

// Accept incomming connections
func (read *SocketReader) StartReader() {
	fmt.Println("Started accepting")
	for {
		conn, err := read.Listener.Accept()
		if err != nil {
			log.Printf("SocketReader could not accept incomming connection: %s", err)
		}

		go read.HandleConnection(conn)
	}
	fmt.Println("No more accepting!")
}

// Handle each incoming connection
func (read *SocketReader) HandleConnection(conn net.Conn) {
	// Read json data from the connection
	var data bytes.Buffer
	var total int
	tmp := make([]byte, 256)

	// Infinitely handle the connection until it's closed
Loop:
	for {
		num, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Printf("SocketReader failed to read: %s", err)
			}
			break
		}
		total += num
		// If the buffer was filled, try to read again to make sure there is no data left
		for num == 256 {
			data.Write(tmp)
			num, err := conn.Read(tmp)
			if err != nil {
				if err != io.EOF {
					log.Printf("SocketReader failed to read: %s", err)
				}
				break Loop
			}
			total += num
		}
		// Append what we have to the end of the message and handle it
		data.Write(tmp)
		read.DispatchConnection(conn, data.Bytes()[:total])

		// Reset the message and message size
		data.Reset()
		total = 0
	}
}

func (read *SocketReader) DispatchConnection(conn net.Conn, apiParams []byte) {
	call := new(JsonRpc)
	err := json.Unmarshal(apiParams, call)
	log.Print(err)
	if err != nil {
		log.Print("Socket reader could not parse command")
		return
	}

	// Extract the json from the call parameters, and encode it as a string
	command, _ := json.Marshal(call.Params)

	// Dispatch function based on api
	if call.Api == "local" {
		result := read.LocalFuncmap[call.Name](read.Local, command)
		var genericResult interface{}
		json.Unmarshal(result.Bytes(), &genericResult)

		resultJson := new(JsonRpcResult)
		resultJson.Name = call.Name
		resultJson.Result = genericResult

		buf, _ := json.Marshal(resultJson)

		conn.Write(buf)
	} else if call.Api == "twitch" {
		result := read.TwitchFuncmap[call.Name](read.Twitch, command)
		var genericResult interface{}
		json.Unmarshal(result.Bytes(), &genericResult)
		// fmt.Println("Stuff:", string(result.Bytes()))
		// fmt.Println("Stuff2:", genericResult)

		resultJson := new(JsonRpcResult)
		resultJson.Name = call.Name
		resultJson.Result = genericResult

		buf, _ := json.Marshal(resultJson)

		conn.Write(buf)
	}
}
