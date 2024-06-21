package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"topiik/internal/command"
	"topiik/internal/proto"
)

func main() {
	// Connect to the server
	//conn, err := net.Dial("tcp", "localhost:8302")
	host := "localhost"
	port := "8302"
	tcpServer, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	// for loop keep cli alive
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimRight(line, " \t\r\n")
		if err != nil {
			break
		}

		strs := strings.SplitN(line, " ", 2)
		if strs[0] == command.QUIT {
			conn.Close()
		}
		// TODO: valid command
		// Send some data to the server
		data, err := proto.Encode(line)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		go response(conn)
	}
}

func response(conn net.Conn) {
	buf := make([]byte, 512)
	/*n, err := conn.Read(buf[0:])
	fmt.Println(n)
	if err != nil {
		fmt.Println(err)
		return
	}*/

	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		if n <= 0 {
			break
		}
		fmt.Printf("%s\n", buf)
	}
}
