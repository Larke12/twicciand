package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
    "github.com/nickvanw/ircx"
    "github.com/sorcix/irc"
)

type Chat interface {
    init()
    setCredentials(user string, channel string, pass string)
}

type wsHandler struct {
    chat *TwitchChat
}

type TwitchChat struct {
    Username *string
    Server   *string
    Channel  *string
    Password *string
    msgsFromClient chan []byte
    msgsToClient chan []byte
}

var channel *string
var upgrader = websocket.Upgrader{
    ReadBufferSize:  2048,
    WriteBufferSize: 2048,
    CheckOrigin: func(r *http.Request) bool {
	return true
    },
}

func (chat *TwitchChat) init() {
    flag.Parse()
    chat.msgsToClient = make(chan []byte, 128)
    chat.msgsFromClient = make(chan []byte, 128)
}

func (chat *TwitchChat) setCredentials(user string, channel string, pass string) {
    chat.Username = flag.String("name", user, "Nick to use in IRC")
    chat.Channel = flag.String("chan", channel, "Channels to join")
    chat.Server = flag.String("server", "irc.twitch.tv:6667", "Host:Port to connect to")
    chat.Password = flag.String("password", "oauth:"+pass, "Connection password")
    fmt.Println("Set chat credentials")
}

func (chat *TwitchChat) startChatServer() {
    bot := ircx.WithLogin(*chat.Server, *chat.Username, *chat.Username, *chat.Password)
    if err := bot.Connect(); err != nil {
	log.Panicln("Unable to dial IRC Server ", err)
    }
    channel = chat.Channel
    chat.RegisterHandlers(bot)
    bot.HandleLoop()
    log.Println("Exiting...")
}

func (chat *TwitchChat) RegisterHandlers(bot *ircx.Bot) {
    bot.HandleFunc(irc.RPL_WELCOME, chat.RegisterConnect)
    bot.HandleFunc(irc.PING, chat.PingHandler)
    bot.HandleFunc(irc.PRIVMSG, chat.PrivMsgHandler)
    fmt.Println("Set chat handlers")
}

func (chat *TwitchChat) RegisterConnect(s ircx.Sender, m *irc.Message) {
    s.Send(&irc.Message{
	Command: irc.JOIN,
	Params:  []string{*channel},
    })
}

func (chat *TwitchChat) PingHandler(s ircx.Sender, m *irc.Message) {
    s.Send(&irc.Message{
	Command:  irc.PONG,
	Params:   m.Params,
	Trailing: m.Trailing,
    })
}

func (chat *TwitchChat) PrivMsgHandler(s ircx.Sender, m *irc.Message) {
    // fmt.Println(m.Prefix.Name + ": " + m.Trailing)
    chat.msgsToClient <- []byte(m.Prefix.Name + ": " + m.Trailing)
}

// Accept incomming connections
func (handle wsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    conn, err := upgrader.Upgrade(w, req, nil)  // omit the responseHeader http.Header for now, not needed
    fmt.Print("Got a connection")
    if err != nil {
	log.Print("Could not open websocket:", err)
    }
    log.Print("Started websocket for chat")
    go handle.chat.SendToClient(conn)
    handle.chat.RecvFromClient(conn)
}

// Write messages from twitch's server to the websocket
func (chat *TwitchChat) SendToClient(conn *websocket.Conn) {
    for msg := range chat.msgsToClient {
	log.Print("Sending to client:", msg)
	err := conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
	    break
	}
    }
    conn.Close()
}

// Write messages from twitch's server to the websocket
func (chat *TwitchChat) RecvFromClient(conn *websocket.Conn) {
    for {
	_, msg, err := conn.ReadMessage()
	log.Print("Received from client:", msg)
	if err != nil {
	    break
	}
	chat.msgsToClient <- msg
    }
    conn.Close()
}
