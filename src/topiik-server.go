//author: duan hongxing
//date: 19 Jun, 2024

package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"topiik/cluster"
	"topiik/executor"
	"topiik/internal/config"
	"topiik/internal/logger"
	"topiik/internal/proto"
	"topiik/node"
	"topiik/persistence"
	"topiik/server"
)

const (
	CONFIG = "--config" // the config file path
)

var serverConfig *config.ServerConfig
var l = logger.Get()

func main() {
	printBanner()
	var err error
	serverConfig, err = readConfig()
	node.SetConfig(*serverConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	// init node
	err = node.InitNode(*serverConfig)
	if err != nil {
		return
	}

	// load controller info on each node
	if node.IsWorker() {
		err = cluster.LoadWorkerGroupInfo()
		if err != nil {
			l.Panic().Msg(err.Error())
		}
		cluster.RequestVote()
	}

	// load data
	if node.IsPersistor() {
		persistence.Load(executor.Execute1)
	}

	// Listen for incoming connections
	ln, err := net.Listen("tcp", serverConfig.Listen)
	if err != nil {
		l.Panic().Msg(err.Error())
	}

	go server.StartServer(serverConfig.Host+":"+serverConfig.Port2, serverConfig)

	// Accept incoming connections and handle them
	l.Info().Msgf("Listen to address %s", serverConfig.Listen)
	for {
		conn, err := ln.Accept()
		if err != nil {
			l.Err(err).Msg(err.Error())
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// close the connection when we're done
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	reader := bufio.NewReader(conn)

	for {
		msg, err := proto.Decode(reader)
		if err != nil {
			if err == io.EOF {
				//fmt.Printf("Client %s connection closed\n", conn.RemoteAddr())
				break
			}
			l.Err(err).Msgf("topiik-server::handleConnection %s", err.Error())
			return
		}
		result := executor.Execute(msg, conn.RemoteAddr().String(), serverConfig)
		conn.Write(result)
	}
}

// Print banner(TODO)
func printBanner() {
	l.Info().Msg("[TOPIIK] Starting Topiik Server...")
}

// Read config values from server.env
func readConfig() (*config.ServerConfig, error) {
	configFile := ""
	if len(os.Args) > 1 {
		if os.Args[1] != CONFIG {
			l.Panic().Msgf("Expect --config, but %s provided\n", os.Args[1])
		}
		if len(os.Args) != 3 {
			l.Panic().Msgf("Expect config file path\n")
		}
		configFile = os.Args[2]
	}

	// Get config
	return config.ParseServerConfig(configFile)
}
