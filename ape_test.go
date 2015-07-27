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
		t.Errorf(
			"name is not \"%s\" - name : \"%s\"",
			name, c.Name())
	}

	if len(c.Args()) != len(args) {
		t.Errorf(
			"args length is not %d - length : %d",
			len(args), len(c.Args()))
	}
	for i, arg := range c.Args() {
		if arg != args[i] {
			t.Errorf(
				"args[%d] is not \"%s\" - args[%d] : \"%s\"",
				i, args[i], i, arg)
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
		t.Errorf(
			"command is not nil - command : %v",
			e.Command())
	}

	if e.messageWithoutName() != strings.TrimSpace(message) {
		t.Errorf(
			"message is not \"%s\" - message : \"%s\"",
			strings.TrimSpace(message), e.messageWithoutName())
	}

	e.buildCommand()
	if e.Command() == nil {
		t.Fatal("command is nil")
	}
	if e.Command().Name() != message1 {
		t.Errorf("command name is not \"%s\" - command name : \"%s\"",
			message1, e.Command().Name())
	}
	if e.Command().Args()[0] != message2 {
		t.Errorf("command args[0] is not \"%s\" - command args[0] : \"%s\"",
			message2, e.Command().Args()[0])
	}
}
