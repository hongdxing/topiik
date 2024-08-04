/*
* author: duan hongxing
* data: 3 Jul 2024
* desc:
*	The interal server run on port: configured port + 10000,
*	The server accept:
*	- Raft rpc
*	- Cluster add node
*	- Worker request for information
 */

package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"topiik/internal/config"
	"topiik/internal/proto"
)

func StartServer(address string, serverConfig *config.ServerConfig) {
	// Listen for incoming connections
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Accept incoming connections and handle them
	fmt.Printf("Listen to address %s\n", address)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn, serverConfig)
	}

}

func handleConnection(conn net.Conn, serverConfig *config.ServerConfig) {
	// Close the connection when we're done
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		msg, err := proto.Decode(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}
		result := Execute(msg[4:], serverConfig)
		conn.Write(result)
	}
}
