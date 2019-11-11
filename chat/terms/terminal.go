package terms

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/shazow/ssh-chat/sshd/terminal"
)

// shell is a container for reading from and writing
// to stdout. This gets passed to a Terminal.
type shell struct {
	r io.Reader
	w io.Writer
}

func (sh *shell) Read(data []byte) (n int, err error) {
	return sh.r.Read(data)
}

func (sh *shell) Write(data []byte) (n int, err error) {
	return sh.w.Write(data)
}

type Terminal struct {
	Nick string
	terminal.Terminal
	oldState *terminal.State
}

// NewTerminal creates a new shazow/ssh-chat terminal that
// can read/write to stdout
func NewTerminal(nick string) *Terminal {
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	sh := &shell{r: os.Stdin, w: os.Stdout}

	term := &Terminal{
		Nick:     nick,
		Terminal: *terminal.NewTerminal(sh, ""),
		oldState: oldState,
	}
	term.Terminal.SetEnterClear(true)
	return term
}

// ListenShell listens on the shell prompt for user input
// and writes the input lines to the connection
func (t *Terminal) ListenShell(conn *websocket.Conn, done chan bool) {
	// defer func() {
	// 	c.conn.Close()
	// }()

	fmt.Println("Ctrl-D to break")
	t.SetPrompt(fmt.Sprintf("[%s]: ", t.Nick))

	// Tell server my nickname
	// Need to feed this back to shell prompt
	msg := Message{Username: t.Nick, Command: "/nick", Message: t.Nick}
	conn.WriteJSON(&msg)

	line, err := t.ReadLine()
	for {
		if err == io.EOF { // Ctrl-D exit so send out the done signal
			t.Write([]byte(line))
			done <- true
			// Wait for done signal with timeout
			// select {
			// case <-c.done:
			// case <-time.After(time.Second):
			// }
			return
		} else if (err != nil && strings.Contains(err.Error(), "control-c break")) || len(line) == 0 {
			line, err = t.ReadLine()
		} else {
			// TODO make command parser helper and resuse
			msg := ParseTerminalMessage(line, t.Nick)
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Println("read:", err)
				return
			}
			line, err = t.ReadLine()
		}
	}
}

func (t *Terminal) WriteShell(buf []byte) error {
	_, err := t.Write([]byte(buf))
	return err
}

// Close releases the terminal from stdin/stdout
func (t *Terminal) Close() {
	fd := int(os.Stdin.Fd())
	terminal.Restore(fd, t.oldState)
}
