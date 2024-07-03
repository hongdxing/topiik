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
	"sync"
	"topiik/internal/config"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/shared"
)

var clusterInitWG sync.WaitGroup
var clusterInitResults = make(map[string]string)

const errClearDataMsg = "err: run FLUSHDB to clear data on node: "

func ClusterInit(addresses []string, serverConfig *config.ServerConfig) (err error) {
	if !IsNodeEmpty() {
		return errors.New(errClearDataMsg + serverConfig.Listen)
	}

	// pre check with each node, whether they can join or not
	for _, addr := range addresses {
		clusterInitWG.Add(1)
		go preCheckPeers(addr)
	}

	// if any address not in the pre check results, then return error
	for _, addr := range addresses {
		if _, ok := clusterInitResults[addr]; !ok {
			return errors.New("err: cannot connect to node: " + addr)
		}
	}

	// if any node pre check failed, then return error
	for addr, result := range clusterInitResults {
		if result == RES_CLUSTER_INIT_FAILED {
			return errors.New(errClearDataMsg + addr)
		} else if result == RES_CLUSTER_INIT_NETWORK_ISSUE {
			return errors.New("err: cannot connect to node: " + addr)
		}
	}
	// if pre check no issue, then start to init cluster
	//err = initPeers(addresses)
	return nil
}

func IsNodeEmpty() bool {
	return len(shared.MemMap) == 0
}

/***
** Confirm other nodes via rpc
**	1) check connectivity
**	2) check whether nodes are empty, i.e. no data
**
**/
func preCheckPeers(address string) {
	defer clusterInitWG.Done()
	conn, err := util.PreapareSocketClient(address)
	defer conn.Close()
	if err != nil {
		clusterInitResults[address] = RES_CLUSTER_INIT_NETWORK_ISSUE
		return
	}

	line := "CLUSTER " + CLUSTER_INIT_PRE_CHECK

	// Enocde
	data, err := proto.Encode(line)
	if err != nil {
		clusterInitResults[address] = RES_CLUSTER_INIT_NETWORK_ISSUE
		fmt.Println(err)
		return
	}

	// Send
	_, err = conn.Write(data)
	if err != nil {
		clusterInitResults[address] = RES_CLUSTER_INIT_NETWORK_ISSUE
		fmt.Println(err)
		return
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		clusterInitResults[address] = RES_CLUSTER_INIT_NETWORK_ISSUE
		return
	} else {
		//fmt.Println(string(buf))
		clusterInitResults[address] = string(buf[:n])
	}
}

/*
func initPeers(addresses []string) error {
	for _, addr := range addresses {

		conn, err := util.PreapareSocketClient(addr)
		defer conn.Close()
		if err != nil {
			return err
		}
		line := "CLUSTER " + CLUSTER_INIT_CONFIRM
		data, _ := proto.Encode(line)
		_, err = conn.Write(data)

		if err != nil {
			fmt.Println(err)
			return err
		}

		buf := make([]byte, 128)
		n, err := conn.Read(buf)
		if err != nil {
			return err
		} else {
			//fmt.Println(string(buf))
		}
	}

}*/
