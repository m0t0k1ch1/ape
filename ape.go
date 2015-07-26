package ape

import (
	"regexp"
	"strings"

	"github.com/thoj/go-ircevent"
)

type callbackFunc func(*Event)

type Event struct {
	*irc.Event
	args []string
}

func (e *Event) Args() []string {
	return e.args
}

type Connection struct {
	*irc.Connection
	channel string
	actions map[string]callbackFunc
}

func (con *Connection) Channel() string {
	return con.channel
}

func (con *Connection) RegisterChannel(channel string) {
	con.channel = channel
}

func (con *Connection) Response(message string) {
	con.Privmsg(con.Channel(), message)
}

func (con *Connection) AddCallback(eventCode string, callback callbackFunc) string {
	return con.Connection.AddCallback(eventCode, func(e *irc.Event) {
		callback(&Event{
			Event: e,
			args:  []string{},
		})
	})
}

func (con *Connection) AddAction(command string, callback callbackFunc) {
	con.actions[command] = callback
}

func (con *Connection) Loop() {
	con.joinChannel()
	con.registerActions()
	con.Connection.Loop()
}

func (con *Connection) joinChannel() string {
	return con.Connection.AddCallback("001", func(e *irc.Event) {
		con.Join(con.Channel())
	})
}

func (con *Connection) registerActions() string {
	return con.AddCallback("PRIVMSG", func(e *Event) {
		// delete own name
		message := regexp.MustCompile(`^(.+: )`).ReplaceAllString(e.Message(), "")

		args := strings.Split(message, " ")
		for command, callback := range con.actions {
			if args[0] == command {
				e.args = args[1:]
				callback(e)
			}
		}
	})
}

func NewConnection(nickname, username string) *Connection {
	return &Connection{
		Connection: irc.IRC(nickname, username),
		actions:    map[string]callbackFunc{},
	}
}
