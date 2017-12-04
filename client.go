package irc

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net"
)

// Connector defines API for connecting to the irc server.
type Connector interface {
	Connect() error
}

// Commander defines the API for sending irc commands to the irc server.
type Commander interface {
	Cmd(name string, params ...string) error
}

// Handler defines the API for registering irc event handlers.
type Handler interface {
	HandleReply(code int, fn HandleFunc)
	HandleCommand(name string, fn HandleFunc)
}

// Client defines the client API for interacting with the irc server.
type Client interface {
	Connector
	Commander
	io.Closer
	io.ReaderFrom
	io.WriterTo
	Handler
}

// client is an implementation of the Client interface.
type client struct {
	host          string
	conn          net.Conn
	replyHandlers map[int]HandleFunc
	cmdHandlers   map[string]HandleFunc
}

// NewClient is the constructor for the client struct.
func NewClient(host string) Client {
	return &client{
		host:          host,
		replyHandlers: map[int]HandleFunc{},
		cmdHandlers:   map[string]HandleFunc{},
	}
}

// Connect attempts to connect to the irc server located at the address specified by the client's host field.
func (c *client) Connect() error {
	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		return err
	}
	c.conn = conn

	return nil
}

// Cmd formats and sends a command to the irc server.
func (c *client) Cmd(name string, params ...string) error {
	buf := bytes.NewBufferString(name)
	for _, v := range params {
		buf.WriteString(" " + v)
	}
	buf.WriteString("\r\n")

	if _, err := c.send(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

// WriteTo reads output from the irc server and prints the output to the Writer passed as an argument (e.g. Stdout).
func (c *client) WriteTo(out io.Writer) (int64, error) {
	if out == nil {
		out = ioutil.Discard
	}

	reader := bufio.NewReader(c.conn)
	count := 0
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return int64(count), err
		}
		count += len(line)

		if err := c.handleEvent(line); err != nil {
			out.Write([]byte(err.Error() + "\r\n"))
			continue
		}
		out.Write(line)
	}
}

// ReadFrom reads input from the Reader passed as an argument (e.g. Stdin) and sends each line to the irc server.
func (c *client) ReadFrom(in io.Reader) (int64, error) {
	reader := bufio.NewReader(in)
	count := 0
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return int64(count), err
		}
		count += len(line)
		if _, err := c.send(append(line, '\r', '\n')); err != nil {
			return int64(count), err
		}
	}
}

// send writes to the connection to the irc server.
func (c *client) send(b []byte) (n int, err error) {
	return c.conn.Write(b)
}

// Close closes the connection to the irc server.
func (c *client) Close() error {
	return c.conn.Close()
}

// HandleFunc is the definition of the irc event handler function.
type HandleFunc func(msg *Event)

// handleEvent parses a stream of bytes sent from the irc server into an Event and then calls any event handlers listening for that command/reply.
func (c *client) handleEvent(raw []byte) error {
	e, err := parseEvent(raw)
	if err != nil {
		return err
	}

	if e.IsCmd() {
		if handler, ok := c.cmdHandlers[e.Cmd]; ok {
			handler(e)
		}
		return nil
	}
	if e.IsReply() {
		if handler, ok := c.replyHandlers[e.Code]; ok {
			handler(e)
		}
	}

	return nil
}

// HandleReply registers an event handler for an irc server reply.
func (c *client) HandleReply(code int, fn HandleFunc) {
	c.replyHandlers[code] = fn
}

// HandleCommand registers an event handler for an irc server command (e.g. PING).
func (c *client) HandleCommand(name string, fn HandleFunc) {
	c.cmdHandlers[name] = fn
}
