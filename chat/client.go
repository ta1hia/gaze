package chat

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	// "golang.org/x/crypto/ssh/terminal"
	"github.com/shazow/ssh-chat/sshd/terminal"
)

var (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = int64(512)
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

// Terminal that reads/writes to stdout
type Client struct {
	nickname string
	terminal.Terminal
	conn     *websocket.Conn
	oldState *terminal.State
	done     chan struct{}
}

// NewClientWithStdInOut creates a new terminal that
// can read/write to stdout
func NewClientWithStdInOut(conn *websocket.Conn, nick string) (cli *Client, err error) {
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	sh := &shell{r: os.Stdin, w: os.Stdout}
	cli = &Client{
		Terminal: *terminal.NewTerminal(sh, ""),
		nickname: nick,
		conn:     conn,
		oldState: oldState,
		done:     make(chan struct{}),
	}
	cli.Terminal.SetEnterClear(true)
	cli.conn.SetReadLimit(maxMessageSize)
	return cli, nil
}

// Close releases the terminal from stdin/stdout
func (c *Client) Close() {
	fd := int(os.Stdin.Fd())
	terminal.Restore(fd, c.oldState)
}

// ListenShell listens on the shell prompt for user input
// and writes the input lines to the connection
func (c *Client) ListenShell() {
	defer func() {
		c.conn.Close()
	}()
	fmt.Println("Ctrl-D to break")
	c.SetPrompt(fmt.Sprintf("[%s]: ", c.nickname))

	line, err := c.ReadLine()
	for {
		if err == io.EOF { // Ctrl-D exit so send out the done signal
			c.Write([]byte(line))
			err := c.conn.Close()
			// err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return
			}

			// Wait for done signal with timeout
			// select {
			// case <-c.done:
			// case <-time.After(time.Second):
			// }
			return
		} else if (err != nil && strings.Contains(err.Error(), "control-c break")) || len(line) == 0 {
			line, err = c.ReadLine()
		} else {
			v := Message{Username: c.nickname, Message: line}
			err := c.conn.WriteJSON(v)
			if err != nil {
				log.Println("read:", err)
				return
			}
			line, err = c.ReadLine()
		}
	}
}

// ListenConnection listens on the connection and writes any
// incoming lines to the terminal
func (c *Client) ListenConnection() {
	// Display msgs
	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			return
		}
		s := fmt.Sprintf("%s: %s\n", msg.Username, msg.Message)
		c.Write([]byte(s))
	}
}
