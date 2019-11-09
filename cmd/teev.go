/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var (
	// hostAddr string
	// channel  string
	// nickname string

	// connectCmd represents the connect command
	teevCmd = &cobra.Command{
		Use:   "teev",
		Short: "connect to a gaze chat server",
		Long: `Gaze provides a client interface for connecting to and 
communicating with a gaze chat server. The client provides a terminal
as a text input/output environment and establishes a websocket client
connection with the gaze server.`,
		RunE: teev,
	}
)

func init() {
	rootCmd.AddCommand(teevCmd)

	// connectCmd.Flags().StringVarP(&hostAddr, "host", "H", "localhost:8844", "Address of the gaze server to connect to")
	// connectCmd.Flags().StringVarP(&channel, "room", "c", "", "Room to connect to; creates the room if it doesn't already exist")
	// connectCmd.Flags().StringVarP(&nickname, "nick", "n", os.Getenv("USER"), "Set nickname")
	// connectCmd.MarkFlagRequired("channel")

}

func teev(cmd *cobra.Command, args []string) error {
	// Initialize application
	app := tview.NewApplication()

	// Create label
	label := tview.NewTextView().SetText("Please enter a nickname: ")

	// Create input field
	input := tview.NewInputField().SetLabel(" ")

	// Create submit button
	btn := tview.NewButton("Submit")

	// Create empty Box to pad each side of appGrid
	bx := tview.NewBox()

	// Create Grid containing the nickname setting widget
	appGrid := tview.NewGrid().
		SetColumns(-1, 24, 16, -1).
		SetRows(-1, 2, 3, -1).
		AddItem(bx, 0, 0, 3, 1, 0, 0, false). // Left - 3 rows
		AddItem(bx, 0, 1, 1, 1, 0, 0, false). // Top - 1 row
		AddItem(bx, 0, 3, 3, 1, 0, 0, false). // Right - 3 rows
		AddItem(bx, 3, 1, 1, 1, 0, 0, false). // Bottom - 1 row
		AddItem(label, 1, 1, 1, 1, 0, 0, false).
		AddItem(input, 1, 2, 1, 1, 0, 0, false).
		AddItem(btn, 2, 1, 1, 2, 0, 0, false)

	// Create text box for chat history
	history := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	// Create text box for connected members
	members := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	// Create input field for chat prompt
	prompt := tview.NewInputField()

	// Create Grid containing the chat's widgets
	chatGrid := tview.NewGrid().
		SetColumns(-10, 0).
		SetRows(-120, 0).
		SetBorders(true).
		AddItem(history, 0, 0, 1, 1, 0, 0, false).
		AddItem(members, 0, 1, 1, 1, 0, 0, false).
		AddItem(prompt, 1, 0, 2, 1, 0, 0, false)

		// submittedName is toggled each time Enter is pressed
	var submittedName bool
	var nick string

	// Capture user input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
					app.SetRoot(chatGrid, true).SetFocus(prompt)
				}

			} else { // In all other cases, handle chat prompt input
				// Clear the input field
				msg := prompt.GetText()
				prompt.SetText("")
				history.Write([]byte(nick + ": " + msg + "\n"))

				// Display appGrid and focus the input field
				app.SetRoot(chatGrid, true).SetFocus(prompt)
			}
			return nil
		case tcell.KeyEsc:
			// Exit the application
			app.Stop()
			return nil
		}

		return event
	})

	// Set the grid as the application root and focus the input field
	app.SetRoot(appGrid, true).SetFocus(input)

	// Run the application
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
