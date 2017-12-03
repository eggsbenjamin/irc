## Examples

### Local Testing

[ngircd](https://ngircd.barton.de) is useful for local testing.

### Simple Client

Run the [simple client](simple_client.go) by entering the following into your terminal from the [examples](../examples) directory.

```
IRC_HOST=${host of the irc server to which you'd like to connect (e.g. localhost:6667)} \
IRC_NICK=${irc nickname} \
IRC_USER=${irc username} \
IRC_PWD=${irc password (optional)}
go run simple_client.go
```
