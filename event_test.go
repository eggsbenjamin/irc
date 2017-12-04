package irc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMessage(t *testing.T) {
	tests := []struct {
		title    string
		input    []byte
		expected *Event
	}{
		{
			"Reply with no params nor content",
			[]byte(`:irc.example.net 001 ben`),
			&Event{
				Sender:    "irc.example.net",
				Recipient: "ben",
				Code:      1,
			},
		},
		{
			"Reply with content but no params",
			[]byte(`:irc.example.net 001 ben :Welcome to the Internet Relay Network ben!~bendavies@localhost`),
			&Event{
				Sender:    "irc.example.net",
				Recipient: "ben",
				Code:      1,
				Content:   []string{"Welcome", "to", "the", "Internet", "Relay", "Network", "ben!~bendavies@localhost"},
			},
		},
		{
			"Reply with params but no content",
			[]byte(`:irc.example.net 004 ben irc.example.net ngircd-24 abBcCFiIoqrRswx abehiIklmMnoOPqQrRstvVz`),
			&Event{
				Sender:    "irc.example.net",
				Recipient: "ben",
				Code:      4,
				Params:    []string{"irc.example.net", "ngircd-24", "abBcCFiIoqrRswx", "abehiIklmMnoOPqQrRstvVz"},
			},
		},
		{
			"Reply with params and content",
			[]byte(":irc.example.net 353 ben = #general :j5ZUwMH745tLJkscthg j5ZUwMH745tLJkscthd j5ZUwMH745tLJkscthm j5ZUwMH745tLJkscthk j5ZUwMH745tLJkscthj j5ZUwMH745tLJkscthi j5ZUwMH745tLJkscthh j5ZUwMH745tLJksc"),
			&Event{
				Sender:    "irc.example.net",
				Recipient: "ben",
				Code:      353,
				Params:    []string{"=", "#general"},
				Content:   []string{"j5ZUwMH745tLJkscthg", "j5ZUwMH745tLJkscthd", "j5ZUwMH745tLJkscthm", "j5ZUwMH745tLJkscthk", "j5ZUwMH745tLJkscthj", "j5ZUwMH745tLJkscthi", "j5ZUwMH745tLJkscthh", "j5ZUwMH745tLJksc"},
			},
		},
		{
			"Cmd with recipient and sender",
			[]byte(`:bob!~bobdavies@localhost PRIVMSG ben :hello`),
			&Event{
				Sender:    "bob!~bobdavies@localhost",
				Recipient: "ben",
				Cmd:       "PRIVMSG",
				Content:   []string{"hello"},
			},
		},
		{
			"Cmd with no recipient nor sender",
			[]byte(`PING :irc.example.net`),
			&Event{
				Cmd:     "PING",
				Content: []string{"irc.example.net"},
			},
		},
	}

	for _, v := range tests {
		t.Run(v.title, func(t *testing.T) {
			actual, err := parseEvent(v.input)

			require.NoError(t, err)
			require.Equal(t, v.expected, actual)
		})
	}
}
