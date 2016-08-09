# ape

[![GoDoc](https://godoc.org/github.com/m0t0k1ch1/ape?status.svg)](https://godoc.org/github.com/m0t0k1ch1/ape) [![wercker status](https://app.wercker.com/status/4de970b70eff735cd6cdd0a1c51d10e1/s/master "wercker status")](https://app.wercker.com/project/bykey/4de970b70eff735cd6cdd0a1c51d10e1)

IRC reaction bot framework based on [thoj/go-ircevent](https://github.com/thoj/go-ircevent)

## Example

``` go
package main

import (
	"log"
	"strings"

	"github.com/m0t0k1ch1/ape"
)

func main() {
	con := ape.NewConnection("nickname", "username")
	con.UseTLS = true
	con.Password = "XXXXX"
	if err := con.Connect("irc.example.com:6667"); err != nil {
		log.Fatal(err)
	}

	con.RegisterChannel("#poyo")

	con.AddAction("piyo", func(e *ape.Event) {
		con.SendMessage("poyo")
	})

	con.AddAction("say", func(e *ape.Event) {
		con.SendMessage(strings.Join(e.Command().Args(), " "))
	})

	con.AddAction("🙏", func(e *ape.Event) {
		con.SendMessage("解脱")
		con.Part(con.Channel())
		con.Join(con.Channel())
		con.SendMessage("輪廻転生")
	})

	con.Loop()
}
```
