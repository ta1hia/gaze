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
	"github.com/spf13/cobra"

	"github.com/tahia-khan/gaze/chat"
)

// serveCmd represents the serve command
var (
	bindAddr string

	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "run a gaze chat server",
		Long:  `Runs a gaze chat server.`,
		Run:   serve,
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&bindAddr, "bind", "b", "localhost:8844", "host:port address to listen on")
}

func serve(cmd *cobra.Command, args []string) {
	server := chat.NewGaze()
	server.Serve(bindAddr)
}
