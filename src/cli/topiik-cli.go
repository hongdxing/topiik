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

const (
	BUF_SIZE = 512
)

/***
**
**
** Syntax: topiik-cli(.exe) --host localhost:8301 [--pass password]
**/
func main() {
	// Connect to the server
	//conn, err := net.Dial("tcp", "localhost:8302")
	host := "localhost:8301"
	pass := ""
	const invalidArgs = "invalid args"
	fmt.Println(os.Args)
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			if strings.ToLower(os.Args[i]) == "--host" {
				if len(os.Args) < i+1 {
					fmt.Println(invalidArgs)
					return
				}
				host = os.Args[i+1]
				i++
			} else if strings.ToLower(os.Args[i]) == "--pass" {
				if len(os.Args) < i+1 {
					fmt.Println(invalidArgs)
					return
				}
				pass = os.Args[i+1]
				i++
				fmt.Println(pass)
			} else {
				fmt.Println(invalidArgs)
			}
		}
	}

	tcpServer, err := net.ResolveTCPAddr("tcp", host)
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
		if strings.ToUpper(strs[0]) == command.QUIT {
			conn.Close()
		}
		// TODO: valid command
		// Enocde
		data, err := proto.Encode(line)
		if err != nil {
			fmt.Println(err)
		}

		// Send
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		go response(conn)
	}
}

func response(conn net.Conn) {
	buf := make([]byte, BUF_SIZE)
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
		fmt.Printf("%s", buf)
		if n <= 0 || n < BUF_SIZE {
			break
		}
	}
	fmt.Println()
}
