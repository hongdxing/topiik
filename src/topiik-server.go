package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"topiik/ccss"
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

	// start CCSS Capital server
	if serverConfig.Role == ccss.CONFIG_ROLE_CAPITAL {
		go ccss.StartServer(serverConfig.Host + ":" + serverConfig.CA_PORT)
	}
	go persistent.Persist(*serverConfig)

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
		/*var result []byte
		if serverConfig.Role == ccss.CONFIG_ROLE_CAPITAL {
			result = ccss.Execute(msg)
		} else {
			result = executer.Execute(msg, serverConfig, nodeStatus)
		}*/
		result := executer.Execute(msg, serverConfig, nodeStatus)
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
	fmt.Printf("Self check start\n")

	var exist bool
	mainPath := getMainPath()
	dataDir := "data"
	slash := string(os.PathSeparator)
	nodeFile := path.Join(mainPath, slash, dataDir, slash, "ccss_node")

	// data dir
	exist, err = util.PathExists(dataDir)
	if err != nil {
		return err
	}

	if !exist {
		fmt.Println("Creating data dir...")
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
	if !exist {
		fmt.Println("Creating node file...")
		var file *os.File
		file, err = os.Create(nodeFile)
		if err != nil {
			//
			return errors.New(res_init_node_failed)
		}
		defer file.Close()

		nodeId = util.RandStringRunes(16)
		if err != nil {
			return errors.New("init node failed")
		}
		fmt.Println(nodeId)
		_, err = file.WriteString(nodeId)
		if err != nil {
			return errors.New("init node failed")
		}
	} else {
		//
	}

	// capital
	if serverConfig.Role == ccss.CONFIG_ROLE_CAPITAL {
		err = initCaptialNode()
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	fmt.Printf("Self check done\n")
	return nil
}

func initCaptialNode() (err error) {
	exist := false                    // whether the file exist
	mainPath := getMainPath()         // path of the server running
	slash := string(os.PathSeparator) // path separator

	var captialMap = make(map[string]ccss.Capital)
	var salorMap = make(map[string]ccss.Salor)
	var partitionMap = make(map[string]ccss.Partition)

	// the capital file
	capitalPath := path.Join(mainPath, slash, DATA_DIR, slash, "ccss_capital")
	exist, err = util.PathExists(capitalPath)
	if err != nil {
		return err
	}
	if !exist {
		fmt.Println("creating capital file...")
		var file *os.File
		file, err = os.Create(capitalPath)
		if err != nil {
			return errors.New(res_init_node_failed)
		}
		defer file.Close()

		capital := ccss.Capital{
			Id:      nodeId,
			Address: serverConfig.Listen,
		}
		captialMap[nodeId] = capital
		var jsonBytes []byte
		jsonBytes, err = json.Marshal(captialMap)
		file.WriteString(string(jsonBytes))
	} else {
		fmt.Println("loading capital metadata...")
		var file *os.File
		file, err = os.Open(capitalPath)
		if err != nil {
			return errors.New(res_init_node_failed)
		}
		defer file.Close()

		captialMap = readMetadata[map[string]ccss.Capital](*file)
		fmt.Println(captialMap)
	}

	// the salor file
	salorPath := path.Join(mainPath, slash, DATA_DIR, slash, "ccss_salor")
	exist, err = util.PathExists(salorPath)
	if err != nil {
		return err
	}
	if !exist {
		fmt.Println("creating salor file...")
		var file *os.File
		file, err = os.Create(salorPath)
		if err != nil {
			return errors.New(res_init_node_failed)
		}
		defer file.Close()
	} else {
		fmt.Println("loading salor metadata...")
		var file *os.File
		file, err = os.Open(salorPath)
		if err != nil {
			return errors.New(res_init_node_failed)
		}
		defer file.Close()

		salorMap = readMetadata[map[string]ccss.Salor](*file)
		fmt.Println(salorMap)
	}

	// the partition file
	patitionPath := path.Join(mainPath, slash, DATA_DIR, slash, "ccss_partition")
	exist, err = util.PathExists(patitionPath)
	if err != nil {
		return err
	}
	if !exist {
		fmt.Println("creating partition file...")
		var file *os.File
		file, err = os.Create(patitionPath)
		if err != nil {
			return errors.New(res_init_node_failed)
		}
		defer file.Close()
	} else {
		fmt.Println("loading partition metadata...")
		var file *os.File
		file, err = os.Open(patitionPath)
		if err != nil {
			return errors.New(res_init_node_failed)
		}
		defer file.Close()

		partitionMap = readMetadata[map[string]ccss.Partition](*file)
		fmt.Println(partitionMap)
	}

	return nil
}

func readMetadata[T any](file os.File) (t T) {
	var jsonBytes = make([]byte, 512)
	var jsonStr string
	for {
		n, err := file.Read(jsonBytes)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		jsonStr += string(jsonBytes[:n])
	}
	if len(jsonStr) > 0 {
		err := json.Unmarshal([]byte(jsonStr), &t)
		if err != nil {
			panic(err)
		}
	}
	return t
}

func getMainPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	mainPath := filepath.Dir(ex)
	return mainPath
}
