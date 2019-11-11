package terms

import (
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

// TVTerminal is a terminal chat ui using tview
type TVTerminal struct {
	app     *tview.Application
	history *tview.TextView
}

// NewTVTerminal creates a new tview terminal for a chat application
func NewTVTerminal() *TVTerminal {
	app := tview.NewApplication() // Initialize application

	history := tview.NewTextView(). // Create text box for chat history
					SetDynamicColors(true).
					SetRegions(true).
					SetScrollable(true).
					SetChangedFunc(func() {
			app.Draw()
		})

	// Create the terminal
	term := &TVTerminal{
		app:     app,
		history: history,
	}

	return term
}

// ListenShell listens on the shell prompt for user input
// and writes the input lines to the connection
func (t *TVTerminal) ListenShell(conn *websocket.Conn, done chan bool) {

	label := tview.NewTextView().SetText("Please enter a nickname: ") // Create label
	input := tview.NewInputField().SetLabel(" ")                      // Create input field
	btn := tview.NewButton("Submit")                                  // Create submit button

	bx := tview.NewBox() // Create empty Box to pad each side of appGrid

	appGrid := tview.NewGrid(). // Create Grid containing the nickname setting widget
					SetColumns(-1, 24, 16, -1).
					SetRows(-1, 2, 3, -1).
					AddItem(bx, 0, 0, 3, 1, 0, 0, false). // Left - 3 rows
					AddItem(bx, 0, 1, 1, 1, 0, 0, false). // Top - 1 row
					AddItem(bx, 0, 3, 3, 1, 0, 0, false). // Right - 3 rows
					AddItem(bx, 3, 1, 1, 1, 0, 0, false). // Bottom - 1 row
					AddItem(label, 1, 1, 1, 1, 0, 0, false).
					AddItem(input, 1, 2, 1, 1, 0, 0, false).
					AddItem(btn, 2, 1, 1, 2, 0, 0, false)

	members := tview.NewTextView(). // Create text box for connected members
					SetDynamicColors(true).
					SetRegions(true).
					SetScrollable(true).
					SetChangedFunc(func() {
			t.app.Draw()
		})

	prompt := tview.NewInputField() // Create input field for chat prompt

	chatGrid := tview.NewGrid(). // Create Grid containing the chat's widgets
					SetColumns(-10, 0).
					SetRows(-120, 0).
					SetBorders(true).
					AddItem(t.history, 0, 0, 1, 1, 0, 0, false).
					AddItem(members, 0, 1, 1, 1, 0, 0, false).
					AddItem(prompt, 1, 0, 2, 1, 0, 0, false)
	// submittedName is toggled each time Enter is pressed
	var submittedName bool
	var nick string

	// Capture user input
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Anything handled here will be executed on the main thread
		switch event.Key() {
		case tcell.KeyEnter:

			// Initial screen - set nickname
			if !submittedName {
				name := input.GetText()
				if strings.TrimSpace(name) == "" { // Stay on this screen until the user enters a valid nickname
					input.SetText("")
				} else { // Go to the chat room!
					submittedName = true
					nick = name
					prompt.SetLabel(fmt.Sprintf("%s: ", name))
					t.app.SetRoot(chatGrid, true).SetFocus(prompt)
				}

			} else { // In all other cases, handle chat prompt input

				// Clear the input field
				line := prompt.GetText()
				prompt.SetText("")

				msg := ParseTerminalMessage(line, nick)
				err := conn.WriteJSON(msg)
				if err != nil {
					log.Println("read:", err)
					return nil
				}

				// Display appGrid and focus the input field
				t.app.SetRoot(chatGrid, true).SetFocus(prompt)
			}
			return nil
		case tcell.KeyEsc:
			// Exit the application
			t.app.Stop()
			return nil
		}

		return event
	})

	// Set the grid as the application root and focus the input field
	t.app.SetRoot(appGrid, true).SetFocus(input)

	// Run the application
	err := t.app.Run()
	if err != nil {
		log.Fatal(err)
	}

}

func (t *TVTerminal) WriteShell(buf []byte) error {
	_, err := t.history.Write([]byte(buf))
	return err
}
