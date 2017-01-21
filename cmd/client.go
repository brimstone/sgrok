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
	"strconv"
	"strings"

	"github.com/brimstone/sgrok/structs"
	"github.com/hashicorp/yamux"
	"github.com/spf13/cobra"
)

func client(cmd *cobra.Command, args []string) {

	// TODO: Work your own magic here
	fmt.Println("client called")

	var st []structs.StreamStruct
	if stdio, _ := cmd.Flags().GetBool("stdio"); stdio {
		st = append(st, structs.StreamStruct{
			Type: structs.StreamTypeSTDIO,
		})
	}

	httpA, err := cmd.Flags().GetStringArray("http")
	for _, http := range httpA {
		httpAddr := strings.SplitN(http, ":", 2)
		httpPort, _ := strconv.Atoi(httpAddr[1])
		st = append(st, structs.StreamStruct{
			Type:        structs.StreamTypeTCP,
			Destination: httpAddr[0],
			Port:        httpPort,
			Direction:   structs.StreamDirectionServer,
		})
	}

	// Get a TCP connection
	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		panic(err)
	}

	// Setup client side of yamux
	session, err := yamux.Client(conn, nil)
	if err != nil {
		panic(err)
	}

	// Open a new stream
	stream, err := session.Open()
	if err != nil {
		panic(err)
	}

	// Stream implements net.Conn
	//stream.Write([]byte("ping"))

	e := gob.NewEncoder(stream)
	err = e.Encode(st)
	if err != nil {
		panic(err)
	}

	for i := range st {
		if st[i].Type == structs.StreamTypeSTDIO {
			go func() { io.Copy(stream, os.Stdin) }()
			go func() { io.Copy(os.Stdout, stream) }()
		}
	}

}

func init() {
	clientCmd := &cobra.Command{
		Use:   "client",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: client,
	}
	RootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	clientCmd.Flags().Bool("listen", false, "Pass in stdio")
	clientCmd.Flags().Bool("stdio", false, "Pass in stdio")
	clientCmd.Flags().StringArray("http", []string{}, "HTTP Connection")

}
