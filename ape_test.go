package ape

import "testing"

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
