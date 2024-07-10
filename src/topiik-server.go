package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"topiik/cluster"
	"topiik/executer"
	"topiik/internal/config"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/persistent"
	"topiik/raft"
)

const (
	BUF_SIZE = 5          // Buffer size that socket read each time
	CONFIG   = "--config" // the config file path
	DATA_DIR = "data"

	res_init_node_failed = "init node failed"
)

var nodeStatus *raft.NodeStatus
var serverConfig *config.ServerConfig
var nodeId string

func main() {
	printBanner()
	var err error
	serverConfig, err = readConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	// self check
	//err = persistent.SelfCheck()
	err = initNode()
	if err != nil {
		return
	}
	nodeStatus = &raft.NodeStatus{Role: raft.ROLE_FOLLOWER, Term: 0}

	// Listen for incoming connections
	ln, err := net.Listen("tcp", serverConfig.Listen)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Start routines
	//go raft.RequestVote(&serverConfig.JoinList, 200, nodeStatus)
	go persistent.Persist(*serverConfig)
	go cluster.StartServer(serverConfig.Host+":"+serverConfig.PORT2, serverConfig)

	// Accept incoming connections and handle them
	fmt.Printf("Listen to address %s\n", serverConfig.Listen)
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
		result := executer.Execute(msg, serverConfig, nodeId, nodeStatus)
		conn.Write(result)
	}
}

/***
** Print banner
**/
func printBanner() {
	fmt.Println("Starting Topiik Server...")
}

/***
** Read config values from server.env
**/
func readConfig() (*config.ServerConfig, error) {
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

func initNode() (err error) {
	fmt.Printf("Topiik: self check start\n")

	var exist bool
	dataDir := "data"
	nodeFile := cluster.GetNodeFilePath()

	// data dir
	exist, err = util.PathExists(dataDir)
	if err != nil {
		return err
	}

	if !exist {
		fmt.Println("creating data dir...")
		err = os.Mkdir(dataDir, os.ModeDir)
		if err != nil {
			fmt.Println(err)
		}
	}

	// node file
	exist, err = util.PathExists(nodeFile)
	if err != nil {
		return err
	}

	var buf []byte
	var node cluster.Node
	if !exist {
		fmt.Println("creating node file...")

		node.Id = util.RandStringRunes(16)
		node.Role = cluster.ROLE_WORKER // new node start as worker by default
		buf, _ = json.Marshal(node)
		err = os.WriteFile(nodeFile, buf, 0644)
		if err != nil {
			panic("loading node failed")
		}
	} else {
		fmt.Println("loading node...")

		buf, err = os.ReadFile(nodeFile)
		if err != nil {
			panic("loading node failed")
		}
		err = json.Unmarshal(buf, &node)
		if err != nil {
			panic("loading node failed")
		}
		fmt.Printf("load node %s\n", node)
	}

	// load controller metadata
	err = cluster.LoadControllerMetadata(&node)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return err
	}
	fmt.Printf("Topiik: self check done\n")
	return nil
}
