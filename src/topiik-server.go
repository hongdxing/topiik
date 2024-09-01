/*
* author: duan hongxing
* date: 19 Jun, 2024
* desc:
 */

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
	if err != nil {
		fmt.Println(err)
		return
	}
	// self check
	//err = persistence.SelfCheck()
	// init node
	err = node.InitNode(*serverConfig)
	if err != nil {
		return
	}

	// load controller info on each node
	err = cluster.LoadControllerInfo()
	if err != nil {
		l.Panic().Msg(err.Error())
	}

	/*
		// online check
		err = online()
		if err != nil {
			l.Panic().Msg(err.Error())
		}*/

	// load metadata
	if node.IsController() {
		err = cluster.LoadMetadata()
		if err != nil {
			l.Panic().Msg(err.Error())
		}
	}

	// Listen for incoming connections
	ln, err := net.Listen("tcp", serverConfig.Listen)
	if err != nil {
		l.Panic().Msg(err.Error())
	}

	// Start routines
	if !node.IsController() {
		persistence.Load(executor.Execute1)
		//go persistence.PersistAsync() // persist
		//go persistence.Sync()    // sync from Partition Leader
	}

	go server.StartServer(serverConfig.Host+":"+serverConfig.Port2, serverConfig)

	// Accept incoming connections and handle them
	l.Info().Msgf("Listen to address %s", serverConfig.Listen)
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

/*
* Print banner(TODO)
 */
func printBanner() {
	l.Info().Msg("[TOPIIK] Starting Topiik Server...")
}

/*
// check if node still valid in the cluster
func online() error {
	// if current node is controller, and the only controller
	if node.GetNodeInfo().Role == node.ROLE_CONTROLLER && len(cluster.GetControllerInfo().Nodes) == 1 {
		return nil
	}
	var connFailCount = 0
	var count = len(cluster.GetControllerInfo().Nodes)
	for _, controller := range cluster.GetControllerInfo().Nodes {
		if node.GetNodeInfo().Id == controller.Id {
			count--
			continue
		}

		err := doOnline(controller.Addr2)
		if err != nil {
			if err.Error() == "CONNECTION" {
				connFailCount++
			}
			continue
		} else {
			return nil
		}
	}
	// if connFailCount ecquals controller count, then this is the first node to start
	if connFailCount == count {
		return nil
	}
	return errors.New(resp.RES_REJECTED)
}

func doOnline(addr2 string) error {
	conn, err := util.PreapareSocketClient(addr2)
	if err != nil {
		l.Warn().Msgf("online connect to controller %s failed", addr2)
		return errors.New("CONNECTION")
	}
	defer conn.Close()

	var buf []byte
	var bbuf = new(bytes.Buffer)
	if err != nil {
		l.Err(err).Msg(err.Error())
	} else {
		binary.Write(bbuf, binary.LittleEndian, consts.RPC_ONLINE)
		buf = append(buf, bbuf.Bytes()...)
		buf = append(buf, []byte(node.GetNodeInfo().ClusterId)...)
	}
	// Enocde
	buf, err = proto.EncodeB(buf)
	if err != nil {
		l.Err(err)
		return err
	}

	// Send
	_, err = conn.Write(buf)
	if err != nil {
		l.Err(err)
		return err
	}

	// Read
	reader := bufio.NewReader(conn)
	buf, err = proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			l.Err(err).Msgf("online %s\n", err)
		}
	}
	if len(buf) < resp.RESPONSE_HEADER_SIZE {
		l.Err(err)
		return err
	}
	rslt := string(buf[resp.RESPONSE_HEADER_SIZE:])
	if rslt == resp.RES_REJECTED {
		l.Warn().Err(errors.New(resp.RES_REJECTED))
		return err
	}
	return nil
}
*/

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
