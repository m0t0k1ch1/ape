package ape

import (
	"regexp"
	"strings"

	"github.com/thoj/go-ircevent"
)

type callbackFunc func(*Event)

type Command struct {
	name string
	args []string
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) Args() []string {
	return c.args
}

func newCommand(name string, args []string) *Command {
	return &Command{
		name: name,
		args: args,
	}
}

type Event struct {
	*irc.Event
	command *Command
}

func (e *Event) Command() *Command {
	return e.command
}

func (e *Event) messageWithoutName() string {
	message := regexp.MustCompile(`^(.+: )`).ReplaceAllString(e.Message(), "")
	return strings.TrimSpace(message)
}

func (e *Event) buildCommand() {
	args := strings.Split(e.messageWithoutName(), " ")
	e.command = newCommand(args[0], args[1:])
}

func newEvent(event *irc.Event) *Event {
	return &Event{
		Event:   event,
		command: nil,
	}
}

type Connection struct {
	*irc.Connection
	channel     string
	actions     map[string]callbackFunc
	initActions []callbackFunc
}

func (con *Connection) Channel() string {
	return con.channel
}

func (con *Connection) RegisterChannel(channel string) {
	con.channel = channel
}

func (con *Connection) SendMessage(message string) {
	con.Privmsg(con.Channel(), message)
}

func (con *Connection) AddCallback(eventCode string, callback callbackFunc) string {
	return con.Connection.AddCallback(eventCode, func(event *irc.Event) {
		callback(newEvent(event))
	})
}

func (con *Connection) AddInitAction(callback callbackFunc) {
	con.initActions = append(con.initActions, callback)
}

func (con *Connection) AddAction(command string, callback callbackFunc) {
	con.actions[command] = callback
}

func (con *Connection) Loop() {
	con.registerInitActions()
	con.registerActions()
	con.Connection.Loop()
}

func (con *Connection) registerInitActions() string {
	return con.AddCallback("001", func(e *Event) {
		con.Join(con.Channel())

		for _, callback := range con.initActions {
			callback(e)
		}
	})
}

func (con *Connection) registerActions() string {
	return con.AddCallback("PRIVMSG", func(e *Event) {
		e.buildCommand()
		for command, callback := range con.actions {
			if e.Command().Name() == command {
				callback(e)
			}
		}
	})
}

func NewConnection(nickname, username string) *Connection {
	return &Connection{
		Connection:  irc.IRC(nickname, username),
		channel:     "",
		initActions: []callbackFunc{},
		actions:     map[string]callbackFunc{},
	}
}
