package irc

import (
	"encoding/json"
	"fmt"
	"log"

	irc "github.com/thoj/go-ircevent"
)

type IrcMessage struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

func Receive(
	server string,
	port int,
	nick string,
	channels []string,
	ch chan<- string,
) {
	conn := irc.IRC(nick, "aun")
	hostAndPort := fmt.Sprintf("%s:%d", server, port)
	conn.VerboseCallbackHandler = true
	conn.Connect(hostAndPort)

	conn.AddCallback("PRIVMSG", func(event *irc.Event) {
		m := IrcMessage{
			From: event.Nick,
			To:   event.Arguments[0],
			Text: event.Message(),
		}
		buf, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		ch <- string(buf)
	})

	for _, channel := range channels {
		conn.Join(channel)
	}
}
