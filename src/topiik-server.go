package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
	"topiik/internal"
	"topiik/internal/command"
	"topiik/internal/proto"
)

const BUF_SIZE = 5

func main() {
	fmt.Println("Starting Topiik Server...")

	// Listen for incoming connections on port 8080
	ln, err := net.Listen("tcp", ":8302")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Listen to port %s\n", "8302")
	// Accept incoming connections and handle them
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
			fmt.Println(err)
			if err == io.EOF {
				break
			}
			return
		}
		fmt.Printf("%s: %s\n", time.Now().Format(internal.DATA_FMT_MICRO_SECONDS), msg)
		conn.Write([]byte(command.OK))
	}

	// Read incoming data
	/*
		buf := make([]byte, BUF_SIZE)
		for {
			n, err := conn.Read(buf[0:])
			fmt.Printf("%s", buf[0:n])
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(err)
				return
			}

			if n <= 0 || n < BUF_SIZE {
				fmt.Println("break")
				break
			}
		}*/

	conn.Write([]byte("Hello from Server你好"))
	process(conn)
}

func process(conn net.Conn) {
	// split into command + arg
	//strs := strings.SplitN(line, " ", 2)
}
