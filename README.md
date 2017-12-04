## IRC

Simple irc client.

See [examples](examples) for full usage.

### Connecting

To connect to an irc server, instantiate a `Client`, passing the irc server address as an argument, and call `Connect.

```
  client := irc.NewClient("localhost:6667")

  if err := client.Connect(); err != nil {
    log.Fatal(err)
  }
  defer client.Close()
```

### Commands

To send a command, call the `CMD` function on a connected `Client.

```
  if err := client.Cmd(irc.NICK, "my_nickname"); err != nil {
    log.Fatal(err)
  }
```
 
### Handlers

There are two types of handler:
- Reply - handles a response to a command sent by the client.
- Command - handles a command sent from the server.

To register a handler, call the respective handler registration function.

```
  client.HandleCommand(irc.PING, func(e *irc.Event) {
    client.Cmd(irc.PONG, "I'm here!")
  })

  client.HandleReply(irc.RPL_WELCOME, func(e *irc.Event) {
    client.Cmd(irc.JOIN, "#channel")
  })
```

__Note__ in order for handlers to be called, the client must be writing server output. See [Writing](#writing).

### Writing

To write output from the server, call the `WriteTo` function passing an `io.Writer` as an argument. 

If `nil` is passed as an argument the server output will be read but not written to an output destination.

``` 
  if err := client.WriteTo(os.Stdout); err != nil {
    log.Fatal(err)
  }
```

### Reading

To read input to send to the server, call the `ReadFrom` function passing an `io.Reader` as an argument.

```
  if err := client.ReadFrom(os.Stdin); err != nil {
    log.Fatal(err)
  }
```

### Constants

All irc response codes and commands, as defined in [RFC 2812](https://tools.ietf.org/html/rfc2812), are exposed by this package.
