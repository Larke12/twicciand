// Reimplements parts of sorcix/ircx make it possible to send messages from a websocket
package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sorcix/irc"
)

type TwitchChat struct {
	channels []*IrcChannel
	curIn    chan []byte
	curOut   chan []byte
}

type IrcChannel struct {
	Name            string
	Reader          *irc.Decoder
	Writer          *irc.Encoder
	Conn            net.Conn
	RawIrcMessages  chan *irc.Message
	PostToChannel   chan []byte
	ReadFromChannel chan []byte
	Config          *IrcConfig
	retries         int
}

type IrcConfig struct {
	Server     string
	Username   string
	Password   string
	MaxRetries int
}

func CreateIrcChannel(name string, cfg *IrcConfig) (*IrcChannel, error) {
	channel := new(IrcChannel)
	channel.Name = name
	channel.RawIrcMessages = make(chan *irc.Message, 128)
	channel.PostToChannel = make(chan []byte, 128)
	channel.ReadFromChannel = make(chan []byte, 128)
	channel.Config = cfg

	err := channel.Connect()
	if err != nil {
		return nil, err
	}
	channel.Reader = irc.NewDecoder(channel.Conn)
	channel.Writer = irc.NewEncoder(channel.Conn)
	err = channel.Login(cfg)
	go channel.RecvLoop()
	go channel.Sort()
	go channel.SendLoop()
	return channel, err
}

func (channel *IrcChannel) Connect() error {
	var err error
	channel.Conn, err = net.Dial("tcp", channel.Config.Server)
	if err != nil {
		return fmt.Errorf("Could not connect to irc server: %s: %s", channel.Config.Server, err)
	}
	return nil
}

func (channel *IrcChannel) Login(cfg *IrcConfig) error {
	log.Print("Logging into channel: ", channel.Name)
	messages := []*irc.Message{}
	// create necessary login messages
	if cfg.Password != "" {
		messages = append(messages, &irc.Message{
			Command: irc.PASS,
			Params:  []string{"oauth:" + cfg.Password},
		})
	}
	messages = append(messages, &irc.Message{
		Command: irc.NICK,
		Params:  []string{cfg.Username},
	})
	messages = append(messages, &irc.Message{
		Command:  irc.USER,
		Params:   []string{cfg.Username, "0", "*"},
		Trailing: cfg.Username,
	})
	// Send login messages
	var err error
	for _, msg := range messages {
		if err = channel.Send(msg); err != nil {
			return err
		}
	}
	return err
}

func (channel *IrcChannel) Send(msg *irc.Message) error {
	err := channel.Writer.Encode(msg)
	return err
}

func (channel *IrcChannel) SendChatMsg(msg string) {
	channel.PostToChannel <- []byte(msg)
}

func (channel *IrcChannel) Reconnect() error {
	if channel.Config.MaxRetries > 0 {
		channel.Conn.Close()
		err := channel.Connect()
		for err != nil && channel.retries < channel.Config.MaxRetries {
			log.Print("Reconnecting channel ", channel.Name)
			duration := time.Duration(math.Pow(2.0, float64(channel.retries))*200) * time.Millisecond
			time.Sleep(duration)
			channel.retries++
		}
		return err
	} else {
		log.Print("Out of retires for channel: ", channel.Name)
		close(channel.RawIrcMessages)
		close(channel.PostToChannel)
		close(channel.ReadFromChannel)
	}
	return nil
}
func (channel *IrcChannel) RecvLoop() {
	for {
		channel.Conn.SetDeadline(time.Now().Add(300 * time.Second))
		msg, err := channel.Reader.Decode()
		if err != nil {
			// TODO: implement reconnect
			channel.Reconnect()
			log.Print("Lost connection to chat channel: ", channel.Name, ": ", err)
			return
		}
		channel.RawIrcMessages <- msg
	}
}

func (channel *IrcChannel) Sort() {
	// Sort and handle irc messages
	for msg := range channel.RawIrcMessages {
		if msg.Command == irc.RPL_WELCOME {
			channel.handleConnect(msg)
		} else if msg.Command == irc.PING {
			channel.handlePing(msg)
		} else if msg.Command == irc.PRIVMSG {
			fmt.Println(msg.Params, ":", msg.Trailing)
			channel.handlePrivMsg(msg)
		}
	}
}

func (channel *IrcChannel) handleConnect(m *irc.Message) {
	channel.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{channel.Name},
	})
}

func (channel *IrcChannel) handlePrivMsg(msg *irc.Message) {
	channel.ReadFromChannel <- []byte(msg.Prefix.Name + ": " + msg.Trailing)
}

func (channel *IrcChannel) handlePing(msg *irc.Message) {
	channel.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   msg.Params,
		Trailing: msg.Trailing,
	})
}

func (channel *IrcChannel) SendLoop() {
	for msg := range channel.PostToChannel {
		channel.Send(&irc.Message{
			Command:  "PRIVMSG",
			Params:   []string{channel.Name},
			Trailing: string(msg),
		})
	}
}

type wsHandler struct {
	chat *TwitchChat
}

// var channel *string
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (chat *TwitchChat) AddChannel(user string, channel string, pass string) *IrcChannel {
	config := &IrcConfig{
		Server:     "irc.twitch.tv:6667",
		Username:   user,
		Password:   pass,
		MaxRetries: 3,
	}
	ircchannel, err := CreateIrcChannel(channel, config)
	if err != nil {
		log.Print("Could not connect to channel: ", channel, ": ", err)
		return nil
	}
	chat.channels = append(chat.channels, ircchannel)
	chat.curIn = ircchannel.PostToChannel
	chat.curOut = ircchannel.ReadFromChannel
	fmt.Println("Added new chat channel")
	return ircchannel
}

// Accept incomming connections
func (handle wsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil) // omit the responseHeader http.Header for now, not needed
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
	for msg := range chat.curOut {
		log.Print("Sending to client:", string(msg))
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
		log.Print("Received from client:", string(msg))
		if err != nil {
			break
		}
		chat.curIn <- msg
		chat.curOut <- msg
	}
	conn.Close()
}
