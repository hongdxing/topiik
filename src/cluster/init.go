/***
** author: duan hongxing
** date: 30 Jun 2024
** desc:
**	return keys
**
**/

package cluster

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"topiik/internal/config"
	"topiik/internal/proto"
	"topiik/shared"
)

var clusterInitWG sync.WaitGroup
var clusterInitResults = make(map[string]string)

const errClearDataMsg = "err: Run FLUSHDB to clear data on node: "

func ClusterInit(addresses []string, serverConfig *config.ServerConfig) (err error) {
	if !IsNodeEmpty() {
		return errors.New(errClearDataMsg + serverConfig.Host)
	}

	for _, addr := range addresses {
		clusterInitWG.Add(1)
		go initPeers(addr)
	}

	for _, addr := range addresses {
		if _, ok := clusterInitResults[addr]; !ok {
			return errors.New("err: cannot connect to node: " + addr)
		}
	}

	fmt.Println(clusterInitResults)

	for addr, result := range clusterInitResults {
		if result == CLUSTER_INIT_FAILED {
			return errors.New(errClearDataMsg + addr)
		} else if result == CLUSTER_INIT_NETWORK_ISSUE {
			return errors.New("err: cannot connect to node: " + addr)
		}
	}
	return nil
}

func IsNodeEmpty() bool {
	return len(shared.MemMap) == 0
}

/***
** Confirm other nodes via rpc
**
**
**
**/
func initPeers(address string) {
	defer clusterInitWG.Done()
	tcpServer, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
	}
	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		//fmt.Println(err)
		clusterInitResults[address] = CLUSTER_INIT_NETWORK_ISSUE
		return
	}
	defer conn.Close()

	line := "CLUSTER __CONFIRM__"

	// Enocde
	data, err := proto.Encode(line)
	if err != nil {
		clusterInitResults[address] = CLUSTER_INIT_NETWORK_ISSUE
		fmt.Println(err)
		return
	}

	// Send
	_, err = conn.Write(data)
	if err != nil {
		clusterInitResults[address] = CLUSTER_INIT_NETWORK_ISSUE
		fmt.Println(err)
		return
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		clusterInitResults[address] = CLUSTER_INIT_NETWORK_ISSUE
		return
	} else {
		//fmt.Println(string(buf))
		clusterInitResults[address] = string(buf[:n])
	}
}
