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
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/tahia-khan/gaze/client"
	"github.com/tahia-khan/gaze/client/terms"
)

var (
	hostAddr string
	channel  string
	nickname string

	// connectCmd represents the connect command
	connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "connect to a gaze chat server",
		Long: `Gaze provides a client interface for connecting to and 
communicating with a gaze chat server. The client provides a terminal
as a text input/output environment and establishes a websocket client
connection with the gaze server.`,
		Run: connect,
	}
)

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().StringVarP(&hostAddr, "host", "H", "localhost:8844", "Address of the gaze server to connect to")
	connectCmd.Flags().StringVarP(&nickname, "nick", "n", os.Getenv("USER"), "Set nickname")

}

func connect(cmd *cobra.Command, args []string) {
	u := url.URL{Scheme: "ws", Host: hostAddr, Path: "connect"}

	// Create the websocket connection
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Create an ssh-chat style terminal ui
	// term := terms.NewTerminal(nickname)
	// defer term.Close()

	// Create a tview terminal ui
	term := terms.NewTVTerminal()

	// Create the gaze client
	cli := client.NewGazeClient(c, term)

	// Run!
	cli.Run()
}
