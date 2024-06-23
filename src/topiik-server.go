package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"topiik/executor"
	"topiik/internal/config"
	"topiik/internal/proto"
	"topiik/raft"
)

const (
	BUF_SIZE = 5          // Buffer size that socket read each time
	CONFIG   = "--config" // the config file path
)

var nodeStatus *raft.NodeStatus
var serverConfig *config.ServerConfig

func main() {
	printBanner()
	serverConfig = readConfig()
	nodeStatus = &raft.NodeStatus{Role: raft.ROLE_FOLLOWER, Term: 0}

	// Listen for incoming connections on port 8080
	ln, err := net.Listen("tcp", serverConfig.Listen)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Listen to address %s\n", serverConfig.Listen)

	go raft.RequestVote(&serverConfig.JoinList, 200, nodeStatus)

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
		cmd, err := proto.Decode(reader)
		if err != nil {
			if err == io.EOF {
				//fmt.Printf("Client %s connection closed\n", conn.RemoteAddr())
				break
			}
			fmt.Println(err)
			return
		}
		//fmt.Printf("%s: %s\n", time.Now().Format(consts.DATA_FMT_MICRO_SECONDS), cmd)
		result := executor.Execute(cmd, serverConfig, nodeStatus)
		conn.Write(result)
	}
}

/***
* Print banner
 */
func printBanner() {
	fmt.Println("Starting Topiik Server...")
}

/***
* Read config values from server.env
 */
func readConfig() *config.ServerConfig {
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
	return config.ParseServerConfig(configFile)
}
