package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"topiik/internal/command"
	"topiik/internal/proto"
	"topiik/resp"
)

var host string
var pass string

/***
**
**
** Syntax: topiik-cli(.exe) --host localhost:8301 [--pass password]
**/
func main() {
	// Connect to the server
	//conn, err := net.Dial("tcp", "localhost:8302")
	host = "localhost:8301"
	pass = ""
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
		println("server not available:", err.Error())
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
		fmt.Print(host + ">")
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
			fmt.Println()
			fmt.Println(err)
		}

		// Send
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println()
			fmt.Println(err)
			return
		}

		response(conn)
	}
}

func response(conn net.Conn) {
	reader := bufio.NewReader(conn)
	buf, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			fmt.Printf("(err): %s\n", err)
		}
		return
	}

	if len(buf) > 4 {
		bufSlice := buf[4:5]
		byteBuf := bytes.NewBuffer(bufSlice)
		var flag int8
		err = binary.Read(byteBuf, binary.LittleEndian, &flag)
		if err != nil {
			fmt.Println("(err):")
		}

		if flag == 1 {
			bufSlice = buf[5:6]
			byteBuf = bytes.NewBuffer(bufSlice)
			var datatype int8
			err = binary.Read(byteBuf, binary.LittleEndian, &datatype)
			if err != nil {
				fmt.Println("(err):")
			}
			if datatype == 1 {
				res := buf[resp.RESPONSE_HEADER_SIZE:]
				fmt.Printf("%s\n", string(res))
			} else if datatype == 2 {
				bufSlice = buf[resp.RESPONSE_HEADER_SIZE:]
				byteBuf := bytes.NewBuffer(bufSlice)
				var result int64
				err = binary.Read(byteBuf, binary.LittleEndian, &result)
				if err != nil {
					fmt.Println("(err):")
				}
				fmt.Printf("%v\n", result)
			} else if datatype == 3 {
				bufSlice = buf[resp.RESPONSE_HEADER_SIZE:]
				//fmt.Println(bufSlice)
				var result []string
				/*byteBuf := bytes.NewBuffer(bufSlice)
				err = binary.Read(byteBuf, binary.LittleEndian, &result)
				if err != nil {
					fmt.Println("(err):")
				}*/

				err = json.Unmarshal(bufSlice, &result)
				if err != nil {
					fmt.Printf("(err):%s\n", err.Error())
				}
				fmt.Printf("%v\n", result)
			} else {
				fmt.Println("(err): invalid response type")
			}

		} else {
			res := buf[resp.RESPONSE_HEADER_SIZE:]
			fmt.Printf("(err):%s\n", res)
		}
	} else {
		fmt.Println("(err): unknown")
	}
}
