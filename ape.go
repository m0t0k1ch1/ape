package ape

import "github.com/thoj/go-ircevent"

type Connection struct {
	*irc.Connection
}

type Event struct {
	*irc.Event
}

func (con *Connection) AddCallback(eventCode string, callback func(*Event)) string {
	return con.Connection.AddCallback(eventCode, func(e *irc.Event) {
		callback(&Event{e})
	})
}

func NewConnection(nickname, username string) *Connection {
	return &Connection{
		irc.IRC(nickname, username),
	}
}
