package irc

import (
	"fmt"
	"strconv"
	"strings"
)

// Event represents an irc event (e.g. command/reply).
type Event struct {
	Sender    string
	Recipient string
	Cmd       string
	Params    []string
	Content   []string
	Code      int
}

// IsCmd indicates if an event is an irc command.
func (r *Event) IsCmd() bool {
	return r.Cmd != ""
}

// IsReply indicates if an event is an irc reply.
func (r *Event) IsReply() bool {
	return r.Code != 0
}

// parseEvent parses a line sent by an irc server into an Event.
func parseEvent(raw []byte) (*Event, error) {
	segments := strings.Split(string(raw), " ")
	if len(segments) < 2 {
		return nil, fmt.Errorf("invalid reply format : '%s'", raw)
	}

	const utf8Colon = 58
	e := &Event{}

	if raw[0] == utf8Colon { // has prefix.
		e.Sender = segments[0][1:]
		segments = segments[1:]
	}

	code, err := strconv.Atoi(segments[0])
	if err != nil {
		e.Cmd = segments[0]
	} else {
		e.Code = code
	}

	remainder := segments[1:]
	for i, v := range remainder {
		if v[0] == utf8Colon { // content starts here.
			remainder[i] = v[1:]
			e.Content = remainder[i:]
			break
		}
		if i == 0 {
			e.Recipient = v
		} else {
			e.Params = append(e.Params, v)
		}
	}

	return e, nil
}
