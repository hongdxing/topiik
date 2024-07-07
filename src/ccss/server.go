/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package ccss

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"topiik/internal/proto"
)

func StartServer(address string) {
	// Listen for incoming connections
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	//captialMap[]

	// start RequestVote routine
	// go RequestVote()

	// Accept incoming connections and handle them
	fmt.Printf("Listen to address %s\n", address)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
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
		result, err := Execute(msg)
		if err != nil {
			fmt.Println(err.Error())
			conn.Write([]byte(err.Error()))
		} else {
			conn.Write(result)
		}
	}
}
