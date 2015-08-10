package ape

import (
	"fmt"
	"strings"
	"testing"
	"time"

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

	if e.targetName() != name {
		t.Errorf(
			"target name is not \"%s\" - target name : \"%s\"",
			name, e.targetName())
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
		t.Errorf(
			"command name is not \"%s\" - command name : \"%s\"",
			message1, e.Command().Name())
	}
	if e.Command().Args()[0] != message2 {
		t.Errorf(
			"command args[0] is not \"%s\" - command args[0] : \"%s\"",
			message2, e.Command().Args()[0])
	}
}

func TestConnection(t *testing.T) {
	channel := "#channel"
	command := "command"

	chanForDone := make(chan string, 1)
	callback1 := func(e *Event) {
		chanForDone <- "callback1"
		t.Log("callback1 is invoked")
	}
	callback2 := func(e *Event) {
		chanForDone <- "callback2"
		t.Log("callback2 is invoked")
	}
	callback3 := func(e *Event) {
		chanForDone <- "callback3"
		t.Log("callback3 is invoked")
	}

	con := NewConnection("nickname", "username")
	if con.Channel() != "" {
		t.Errorf(
			"channel is not \"\" - channel : \"%s\"",
			con.Channel())
	}
	if len(con.initActions) > 0 {
		t.Errorf(
			"initActions length is over 0 - initActions length : %d",
			len(con.initActions))
	}
	if con.defaultAction != nil {
		t.Errorf(
			"defaultAction is not nil - defaultAction : %v",
			con.defaultAction)
	}
	if len(con.actions) > 0 {
		t.Errorf(
			"actions length is over 0 - actions length : %d",
			len(con.actions))
	}

	con.RegisterChannel(channel)
	if con.Channel() != channel {
		t.Errorf(
			"channel is not \"%s\" - channel : \"%s\"",
			channel, con.Channel())
	}

	con.AddInitAction(callback1)
	if len(con.initActions) != 1 {
		t.Errorf(
			"initActions length is not 1 - initActions length : %d",
			len(con.initActions))
	}
	con.initActions[0](&Event{})
	if result := <-chanForDone; result != "callback1" {
		t.Errorf(
			"result is not \"callback1\" - result : \"%s\"",
			result)
	}

	con.AddDefaultAction(callback2)
	if con.defaultAction == nil {
		t.Errorf("defaultAction is nil")
	}
	con.defaultAction(&Event{})
	if result := <-chanForDone; result != "callback2" {
		t.Errorf(
			"result is not \"callback2\" - result : \"%s\"",
			result)
	}

	con.AddAction(command, callback3)
	if len(con.actions) != 1 {
		t.Errorf(
			"actions length is not 1 - actions length : %d",
			len(con.actions))
	}
	con.actions[command](&Event{})
	if result := <-chanForDone; result != "callback3" {
		t.Errorf(
			"result is not \"callback3\" - result : \"%s\"",
			result)
	}
}

func TestAction(t *testing.T) {
	server := "irc.freenode.net:6667"
	channel := "#m0t0k1ch1-ape"

	name1 := "ape1"
	name2 := "ape2"

	isReadyCon2 := false
	isInvokeDefaultAction := false
	count := 0

	chanForInit := make(chan bool, 1)
	chanForDefaultAction := make(chan bool, 1)
	chanForCountUp := make(chan bool, 1)
	chanForDone := make(chan bool, 1)

	con1 := NewConnection(name1, name1)
	if err := con1.Connect(server); err != nil {
		t.Fatal(err)
	}
	con1.RegisterChannel(channel)
	con1.AddCallback("366", func(event *Event) {
		ticker := time.NewTicker(1 * time.Second)
		i := 5

		for {
			<-ticker.C
			if isReadyCon2 {
				con1.SendMessage(fmt.Sprintf("%s: count-up", name2))
				i--
				t.Log("con1 - send message to con2")
			}
			if i == 0 {
				ticker.Stop()
				con1.SendMessage(fmt.Sprintf("%s: poyo", name2))
				con1.SendMessage(fmt.Sprintf("%s: quit", name2))
				con1.Quit()
				t.Log("con1 - quit")
			}
		}
	})

	con2 := NewConnection(name2, name2)
	if err := con2.Connect(server); err != nil {
		t.Fatal(err)
	}
	con2.RegisterChannel(channel)
	con2.AddInitAction(func(event *Event) {
		chanForInit <- true
		t.Log("con2 - init")
	})
	con2.AddDefaultAction(func(event *Event) {
		chanForDefaultAction <- true
		t.Log("con2 - default action")
	})
	con2.AddAction("count-up", func(e *Event) {
		chanForCountUp <- true
		t.Log("con2 - count up")
	})
	con2.AddAction("quit", func(e *Event) {
		con2.Quit()
		chanForDone <- true
		t.Log("con2 - quit")
	})

	go con1.Loop()
	go con2.Loop()

	func() {
		for {
			select {
			case <-chanForInit:
				isReadyCon2 = true
			case <-chanForDefaultAction:
				isInvokeDefaultAction = true
			case <-chanForCountUp:
				count++
			case <-chanForDone:
				return
			}
		}
	}()

	if !isInvokeDefaultAction {
		t.Errorf("default action is not invoked")
	}
	if count != 5 {
		t.Errorf("count is not 5 - count : %d", count)
	}
}
