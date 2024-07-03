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

var nodeStatus *NodeStatus
var salorAddress *[]string

func StartServer(address string) {
	// Listen for incoming connections
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// start RequestVote routine
	if nodeStatus.Role == CCSS_ROLE_CO {
		go RequestVote()
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
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	// Close the connection when we're done
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		//cmd, err := proto.Decode(reader)
		msg, err := proto.Decode(reader)
		if err != nil {
			if err == io.EOF {
				//fmt.Printf("Client %s connection closed\n", conn.RemoteAddr())
				break
			}
			fmt.Println(err)
			return
		}
		//fmt.Printf("%s: %s\n", time.Now().Format(consts.DATA_FMT_MICRO_SECONDS), cmd)
		result := execute(msg)
		conn.Write(result)
	}
}
