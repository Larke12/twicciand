// Reimplements parts of sorcix/ircx make it possible to send messages from a websocket
package main

import (
	//	"bytes"
	"fmt"
	"html"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"time"
	"regexp"

	"github.com/gorilla/websocket"
	"github.com/sorcix/irc"			// IRC v3 branch
)

type TwitchChat struct {
	channels	[]*IrcChannel
	auth		*TwitchAuth
	curIn		chan []byte
	curOut		chan []byte
	colorMap	map[string]string
	mod		[]string	// 0 or 1, unused right now
	turbo		[]string	// 0 or 1
	sub		[]string	// 0 or 1
	usertype	[]string	// empty, mod, global_mod, admin or staff
	disp_name	[]string	// users sylized name
	color		[]string	// empty or hexadecimal
	raw		string		// temp string to hold the original msg
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
	messages := []*irc.Message{}
	//log.Print("Logging into channel: ", channel.Name)
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
		log.Print("Out of retries for channel: ", channel.Name)
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
			channel.handleCAP(msg)
			channel.handleConnect(msg)
		} else if msg.Command == irc.PING {
			channel.handlePing(msg)
		} else if msg.Command == irc.PRIVMSG {
			//fmt.Println(msg.Params, ":", msg.Trailing)
			channel.handlePrivMsg(msg)
		}
		// Catch USERSTATE on join
	}
}

func (channel *IrcChannel) handleConnect(m *irc.Message) {
	channel.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{channel.Name},
	})
}

func (channel *IrcChannel) handleCAP(m *irc.Message) {
	// Registering for IRCv3 Membership
	channel.Send(&irc.Message{
		Command: irc.CAP,
		Params: []string{irc.CAP_REQ + " twitch.tv/membership"},
	})
	// Registering for IRCv3 Tags
	channel.Send(&irc.Message{
		Command: irc.CAP,
		Params: []string{irc.CAP_REQ + " twitch.tv/tags"},
	})
	// Registering for IRCv3 Commands
	channel.Send(&irc.Message{
		Command: irc.CAP,
		Params: []string{irc.CAP_REQ + " twitch.tv/commands"},
	})
}

func (channel *IrcChannel) handlePrivMsg(msg *irc.Message) {
	//fmt.Println(msg)
	fmt_msg := new(TwitchChat)
	fmt_msg.raw = msg.String()

	// Parse the tags out of the PRIVMSG for use in the front end

	// Parse usertype
	reUserType, err := regexp.Compile(`user-type\=(.*?)(\;|\s)`)
	if err != nil {
		log.Print("Could not parse UserType\n")
	}
	fmt_msg.usertype = reUserType.FindStringSubmatch(fmt_msg.raw)

	// Parse subscriber
	reSub, err := regexp.Compile(`subscriber\=(.*?)(\;|\s)`)
	if err != nil {
		log.Print("Could not parse Subscriber\n")
	}
	fmt_msg.sub = reSub.FindStringSubmatch(fmt_msg.raw)

	// Parse turbo
	reTurbo, err := regexp.Compile(`turbo\=(.*?)(\;|\s)`)
	if err != nil {
		log.Print("Could not parse Turbo\n")
	}
	fmt_msg.turbo = reTurbo.FindStringSubmatch(fmt_msg.raw)

	// Parse display name
	reDisp, err := regexp.Compile(`display-name\=(.*?)(\;|\s)`)
	if err != nil {
		log.Print("Could not parse DisplayName\n")
	}
	fmt_msg.disp_name = reDisp.FindStringSubmatch(fmt_msg.raw)

	// Parse color tag
	reColor, err := regexp.Compile(`#[[:xdigit:]]{6}`)
	if err != nil {
		log.Print("Could not parse Color\n")
	}
	fmt_msg.color = reColor.FindStringSubmatch(fmt_msg.raw)

	if len(fmt_msg.color) == 1 && len(fmt_msg.disp_name) >= 1 && len(fmt_msg.sub) >= 1 && len(fmt_msg.turbo) >= 1 && len(fmt_msg.usertype) >= 1 {
		// User has all fields (mod or staff)
		// <a href='https://www.twitch.tv/" + msg.Prefix.Name + "/profile' target='_blank'><strong>" + fmt_msg.disp_name[1] + "</strong></a>
		channel.ReadFromChannel <- []byte("<span data-usertype='" + fmt_msg.usertype[1] + "' data-sub='" + fmt_msg.sub[1] + "' data-turbo='" + fmt_msg.turbo[1] + 
		"' style='color:" + fmt_msg.color[0] + "' id='username'><strong>" + fmt_msg.disp_name[1] + "</strong></span><span id='text'>: " + html.EscapeString(msg.Trailing) + " </span>")
	} else if len(fmt_msg.color) == 1 && len(fmt_msg.disp_name) >= 1 && len(fmt_msg.sub) >= 1 && len(fmt_msg.turbo) >= 1 {
		// User is missing user-type tag (non-mod)
		channel.ReadFromChannel <- []byte("<span data-sub='" + fmt_msg.sub[1] + "' data-turbo='" + fmt_msg.turbo[1] + "' style='color:" + fmt_msg.color[0] + 
		"' id='username'><strong>" + fmt_msg.disp_name[1] + "</strong></span><span id='text'>: " + html.EscapeString(msg.Trailing) + " </span>")
	} else if len(fmt_msg.color) == 1 && len(fmt_msg.disp_name) >= 1 {
		// User is missing user-type, subscriber, and turbo tags (rare)
		channel.ReadFromChannel <- []byte("<span data-sub='0' data-turbo='0' style='color:" + fmt_msg.color[0] + 
		"' id='username'><strong>" + fmt_msg.disp_name[1] + "</strong></span><span id='text'>: " + html.EscapeString(msg.Trailing) + " </span>")
	} else if len(fmt_msg.color) == 1 {
		// User is bot (or not authenticated)
		channel.ReadFromChannel <- []byte("<span style='color:" + fmt_msg.color[0] + "' id='username'><strong>" + msg.Prefix.Name + "</strong></span><span id='text'>: " + html.EscapeString(msg.Trailing) + " </span>")
	} else {
		// Randomize colors if the user has never set them before
		rand.Seed(time.Now().UTC().UnixNano())
		colors := []string{
			"#FF0000",
			"#0000FF",
			"#008000",
			"#B22222",
			"#FF7F50",
			"#9ACD32",
			"#FF4500",
			"#2E8B57",
			"#DAA520",
			"#D2691E",
			"#5F9EA0",
			"#1E90FF",
			"#FF69B4",
			"#8A2BE2",
			"#00FF7F",
		}

		/* Map colors to the name, broken for some reason
		color, ok := fmt_msg.colorMap[msg.Prefix.Name]
		if !ok {
			color = colors[rand.Intn(len(colors))]
			fmt_msg.colorMap[msg.Prefix.Name] = color
		}*/

		color := colors[rand.Intn(len(colors))]

		if len(fmt_msg.disp_name) >= 1 && len(fmt_msg.sub) >= 1 && len(fmt_msg.turbo) >= 1 && len(fmt_msg.usertype) >= 1 {
			// User has all fields (mod or staff)
			channel.ReadFromChannel <- []byte("<span data-usertype='" + fmt_msg.usertype[1] + "' data-sub='" + fmt_msg.sub[1] + "' data-turbo='" + fmt_msg.turbo[1] + 
			"' style='color:" + color + "' id='username'><strong>" + fmt_msg.disp_name[1] + "</strong></span><span id='text'>: " + html.EscapeString(msg.Trailing) + " </span>")
		} else if len(fmt_msg.disp_name) >= 1 && len(fmt_msg.sub) >= 1 && len(fmt_msg.turbo) >= 1 {
			// User is missing user-type tag (non-mod)
			channel.ReadFromChannel <- []byte("<span data-sub='" + fmt_msg.sub[1] + "' data-turbo='" + fmt_msg.turbo[1] + "' style='color:" + color + 
			"' id='username'><strong>" + fmt_msg.disp_name[1] + "</strong></span><span id='text'>: " + html.EscapeString(msg.Trailing) + " </span>")
		} else if len(fmt_msg.disp_name) >= 1 {
			// User is missing user-type, subscriber, and turbo tags (rare)
			channel.ReadFromChannel <- []byte("<span data-sub='0' data-turbo='0' style='color:" + color + 
			"' id='username'><strong>" + fmt_msg.disp_name[1] + "</strong></span><span id='text'>: " + html.EscapeString(msg.Trailing) + " </span>")
		}
	}
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

func (channel *IrcChannel) Disconnect() {
	close(channel.PostToChannel)
	close(channel.ReadFromChannel)
	close(channel.RawIrcMessages)
	channel.Conn.Close()
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
		Server:	    "irc.chat.twitch.tv:80",
		Username:   user,
		Password:   pass,
		MaxRetries: 3,
	}
	ircchannel, err := CreateIrcChannel(channel, config)
	if err != nil {
		log.Print("Could not connect to channel: ", channel, ": ", err)
		return nil
	}
	for _, oldchan := range chat.channels {
		oldchan.Disconnect()
	}
	chat.channels = chat.channels[:0]
	chat.channels = append(chat.channels, ircchannel)
	chat.curIn = ircchannel.PostToChannel
	chat.curOut = ircchannel.ReadFromChannel
	//fmt.Println("Added new chat channel")
	return ircchannel
}

// Accept incomming connections
func (handle wsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil) // omit the responseHeader http.Header for now, not needed
	//fmt.Print("Got a connection\n")
	if err != nil {
		log.Print("Could not open websocket:", err)
	}
	//log.Print("Started websocket for chat")
	go handle.chat.SendToClient(conn)
	handle.chat.RecvFromClient(conn)
}

// Write messages from twitch's server to the websocket
func (chat *TwitchChat) SendToClient(conn *websocket.Conn) {
	for msg := range chat.curOut {
		//log.Print("Sending to client: ", string(msg))
		err := conn.WriteMessage(websocket.TextMessage, []byte(string(msg)))

		if err != nil {
			break
		}
	}
	conn.Close()
}

// Write messages from twitch's server to the websocket
/*func (chat *TwitchChat) RecvFromClient(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		//log.Print("Received from client:", string(msg))
		if err != nil {
			break
		}

		//Randomize colors if the user has never set them before
		rand.Seed(time.Now().UTC().UnixNano())
		colors := []string{
			"#FF0000",
			"#0000FF",
			"#008000",
			"#B22222",
			"#FF7F50",
			"#9ACD32",
			"#FF4500",
			"#2E8B57",
			"#DAA520",
			"#D2691E",
			"#5F9EA0",
			"#1E90FF",
			"#FF69B4",
			"#8A2BE2",
			"#00FF7F",
		}

		color := colors[rand.Intn(len(colors))]

		arr := []byte("<span data-sub='0' data-turbo='0' style='color:" + color + "'id='username'><strong>" + chat.auth.Username + ": ")
		msg = []byte("</strong></span><span id='text'>: " + string(msg) + " </span>")

		chat.curIn <- msg
		chat.curOut <- []byte(string((append(arr[:], msg[:]...))))
	}
	conn.Close()
}*/
func (chat *TwitchChat) RecvFromClient(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		log.Print("Received from client:", string(msg))
		if err != nil {
			break
		}
		arr := []byte(chat.auth.Username + ": ")
		chat.curIn <- msg
		chat.curOut <- []byte(string((append(arr[:], msg[:]...))))
	}
	conn.Close()
}

