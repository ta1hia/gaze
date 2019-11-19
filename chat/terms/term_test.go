package terms

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/shazow/ssh-chat/sshd/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/tahia-khan/gaze/chat/terms"
)

// TermSuite runs a suite of tests against a TerminalUI implementation
func TermSuite(t *testing.T, newTerm func() terms.TerminalUI) {

	t.Run("WriteBytes", func(t *testing.T) {
		s := []byte("hello hello")
		term := newTerm()
		err := term.WriteShell(s)
		assert.Nil(t, err, "%+v", err)
	})
}

// TestTerminal runs TermSuite against the Terminal
func TestTerminal(t *testing.T) {
	t.Run("Terminal", func(t *testing.T) {
		TermSuite(t, func() terms.TerminalUI {
			buf := bufio.NewReadWriter(bufio.NewReader(bytes.NewBuffer(nil)), bufio.NewWriter(bytes.NewBuffer(nil)))
			return &terms.Terminal{
				Nick:     "somenick",
				Terminal: *terminal.NewTerminal(buf, ""),
			}
		})
	})
}

// TestTerminal runs TermSuite against the TVTerminal
func TestTVTerminal(t *testing.T) {
	t.Run("TVTerminal", func(t *testing.T) {
		TermSuite(t, func() terms.TerminalUI {
			return terms.NewTVTerminal()
		})
	})
}

func TestParseTerminalMessage(t *testing.T) {
	tests := map[string]struct {
		input string
		nick  string
		want  *terms.Message
	}{
		"message":           {input: "just a regular message", nick: "bob", want: &terms.Message{"bob", "", "just a regular message"}},
		"command no args":   {input: "/yer", nick: "bob", want: &terms.Message{"bob", "/yer", ""}},
		"command with args": {input: "/cmd a command", nick: "bob", want: &terms.Message{"bob", "/cmd", "a command"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := terms.ParseTerminalMessage(tc.input, tc.nick)
			assert.Equal(t, tc.want, got)
		})
	}
}
