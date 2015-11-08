package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type Chat interface {
	init()
	setCredentials(user string, channel string, pass string)
}

type TwitchChat struct {
	Username *string
	Server   *string
	Channel  *string
	Password *string
}

var channel *string

func (chat *TwitchChat) init() {
	flag.Parse()
}

func (chat *TwitchChat) setCredentials(user string, channel string, pass string) {
	chat.Username = flag.String("name", user, "Nick to use in IRC")
	chat.Channel = flag.String("chan", channel, "Channels to join")
	chat.Server = flag.String("server", "irc.twitch.tv:6667", "Host:Port to connect to")
	chat.Password = flag.String("password", "oauth:"+pass, "Connection password")
	fmt.Println("Set credentials")
}

func (chat *TwitchChat) startChatServer() {
	bot := ircx.WithLogin(*chat.Server, *chat.Username, *chat.Username, *chat.Password)
	if err := bot.Connect(); err != nil {
		log.Panicln("Unable to dial IRC Server ", err)
	}
	channel = chat.Channel
	RegisterHandlers(bot)
	bot.HandleLoop()
	log.Println("Exiting..")
}

func RegisterHandlers(bot *ircx.Bot) {
	bot.HandleFunc(irc.RPL_WELCOME, RegisterConnect)
	bot.HandleFunc(irc.PING, PingHandler)
	bot.HandleFunc(irc.PRIVMSG, PrivMsgHandler)
	fmt.Println("Set handlers")
}

func RegisterConnect(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{*channel},
	})
}

func PingHandler(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   m.Params,
		Trailing: m.Trailing,
	})
}

func PrivMsgHandler(s ircx.Sender, m *irc.Message) {
	fmt.Println(m.Prefix.Name + ": " + m.Trailing)
}
