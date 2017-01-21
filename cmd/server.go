// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/brimstone/sgrok/structs"
	"github.com/hashicorp/yamux"
	"github.com/spf13/cobra"
)

//var []structs.StreamStruct

func handleSession(session *yamux.Session) {
	// Accept a stream
	stream, err := session.Accept()
	if err != nil {
		panic(err)
	}

	// Listen for a message
	//buf := make([]byte, 4)
	//stream.Read(buf)
	//fmt.Println("Read", string(buf))
	var st []structs.StreamStruct
	e := gob.NewDecoder(stream)
	err = e.Decode(&st)
	if err != nil {
		panic(err)
	}

	fmt.Println("Client requests", len(st), "streams")

	for i := range st {
		if st[i].Type == structs.StreamTypeSTDIO {
			go func() { io.Copy(stream, os.Stdin) }()
			go func() { io.Copy(os.Stdout, stream) }()
			continue
		}
		//fmt.Println(st[0].Destination)
	}

}

func server(cmd *cobra.Command, args []string) {
	// TODO: Work your own magic here
	fmt.Println("server called")
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		panic(err)
	}

	for {
		// Accept a TCP connection
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		// Setup server side of yamux
		session, err := yamux.Server(conn, nil)
		if err != nil {
			panic(err)
		}

		go handleSession(session)
	}

}

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: server,
	})

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
