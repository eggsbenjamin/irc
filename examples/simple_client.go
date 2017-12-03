package main

import (
	"log"
	"os"

	"github.com/eggsbenjamin/irc"
)

func main() {
	// create client.
	client := irc.NewClient(os.Getenv("IRC_HOST"))

	// connect to server.
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// create handler to respond to irc server 'PING' command.
	client.HandleCommand(irc.PING, func(e *irc.Event) {
		client.Cmd(irc.PONG, "I'm here!")
	})

	// set nickname and username.
	nick, user := os.Getenv("IRC_NICK"), os.Getenv("IRC_USER")
	if err := client.Cmd(irc.NICK, nick); err != nil {
		log.Fatal(err)
	}
	if err := client.Cmd(irc.USER, user, "0 *:", user); err != nil {
		log.Fatal(err)
	}

	// send messages from stdin to irc server.
	go func() {
		log.Fatal(client.ReadFrom(os.Stdin))
	}()

	// write messages from irc server to stdout.
	log.Fatal(client.WriteTo(os.Stdout))
}
