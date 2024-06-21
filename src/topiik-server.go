package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
	"topiik/internal/command"
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/internal/proto"
)

const (
	BUF_SIZE = 5          // Buffer size that socket read each time
	CONFIG   = "--config" // the config file path
)

func main() {
	printMark()
	configFile := ""
	if len(os.Args) > 1 {
		if os.Args[1] != CONFIG {
			fmt.Printf("Expect --config, but %s provided\n", os.Args[1])
			os.Exit(1)
		}
		if len(os.Args) != 3 {
			fmt.Printf("Expect config file path\n")
			os.Exit(1)
		}
		configFile = os.Args[2]
	}

	// Get config
	serverConfig := config.ParseServerConfig(configFile)

	// Listen for incoming connections on port 8080
	ln, err := net.Listen("tcp", serverConfig.Host+":"+serverConfig.Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Listen to port %s\n", serverConfig.Port)
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
		fmt.Printf("%s: %s\n", time.Now().Format(consts.DATA_FMT_MICRO_SECONDS), msg)
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

func printMark() {
	fmt.Println("Starting Topiik Server...")
}
