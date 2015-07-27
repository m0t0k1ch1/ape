package ape

import (
	"fmt"
	"strings"
	"testing"

	"github.com/thoj/go-ircevent"
)

func TestCommand(t *testing.T) {
	name := "name"
	args := []string{"arg1", "arg2"}
	c := newCommand(name, args)

	if c.Name() != name {
		t.Errorf("invalid name - name: \"%s\"", c.Name())
	}

	if len(c.Args()) != len(args) {
		t.Errorf("invalid args length - length: %d", len(c.Args()))
	}
	for i, arg := range c.Args() {
		if arg != args[i] {
			t.Errorf("invalid arg - arg: \"%s\"", arg)
		}
	}
}

func TestEvent(t *testing.T) {
	name := "name"
	message1 := "message1"
	message2 := "message2"

	message := fmt.Sprintf("  %s %s  ", message1, message2)

	e := newEvent(&irc.Event{
		Arguments: []string{fmt.Sprintf("%s: %s", name, message)},
	})
	if e.Command() != nil {
		t.Errorf("command is not nil - command: %v", e.Command())
	}

	if e.messageWithoutName() != strings.TrimSpace(message) {
		t.Errorf("invalid message - message: \"%s\"", e.messageWithoutName())
	}

	e.buildCommand()
	if e.Command() == nil {
		t.Fatal("no command")
	}
	if e.Command().Name() != message1 {
		t.Errorf("invalid command name - name: \"%s\"", e.Command().Name())
	}
	if e.Command().Args()[0] != message2 {
		t.Errorf("invalid command arg - arg: \"%s\"", e.Command().Args()[0])
	}
}
