/***
**
**
**
**
**/

package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"topiik/internal/config"
	"topiik/internal/util"
)

const cluster_init_failed = "cluster init failed"

func ClusterInit(partitions uint16, replicas uint16, serverConfig *config.ServerConfig) (err error) {
	fmt.Println("ClusterInit start...")

	// 0. init cluster
	initCluster(partitions, replicas, serverConfig)

	// 1. open node file
	nodePath := GetNodeFilePath()

	// 2. read from node file
	var jsonStr string
	var buf []byte

	buf, err = os.ReadFile(nodePath)
	if err != nil {
		return errors.New(cluster_init_failed)
	}
	jsonStr = string(buf)

	// 3. unmarshal
	err = json.Unmarshal([]byte(jsonStr), &nodeInfo)
	if err != nil {
		return errors.New(cluster_init_failed)
	}
	if len(nodeInfo.ClusterId) > 0 { // check if current node already in cluster or not
		return errors.New("current node already in a cluster:" + nodeInfo.ClusterId)
	}
	nodeInfo.ClusterId = clusterInfo.Id

	// 4. marshal
	buf2, _ := json.Marshal(nodeInfo)

	// 5. write back to file
	os.Truncate(nodePath, 0)
	os.WriteFile(nodePath, buf2, 0644)

	nodeStatus.Role = RAFT_LEADER
	err = initControllerNode(nodeInfo.Id, serverConfig)
	if err != nil {
		return err
	}
	// after init, the node default is LEADER, and start to AppendEntries()
	go AppendEntries()
	fmt.Println("ClusterInit end")
	return nil
}

func initCluster(partitions uint16, replicas uint16, serverConfig *config.ServerConfig) error {
	if len(clusterInfo.Id) > 0 {
		return errors.New("current node already in cluster:" + clusterInfo.Id)
	}
	// set clusterInfo
	clusterInfo.Id = util.RandStringRunes(16)
	clusterInfo.Ptns = partitions
	clusterInfo.Rpls = replicas

	addrSplit, err := util.SplitAddress(serverConfig.Listen)
	if err != nil {
		panic(err) // if cannot resovle the Address, it's severe error
	}
	clusterInfo.Ctls[nodeInfo.Id] = NodeSlim{
		Id:    nodeInfo.Id,
		Addr:  serverConfig.Listen,
		Addr2: addrSplit[0] + ":" + addrSplit[2]}

	// persist cluster metadata
	clusterPath := GetClusterFilePath()
	_ = os.Remove(clusterPath) // just incase
	data, err := json.Marshal(clusterInfo)
	if err != nil {
		return errors.New(cluster_init_failed)
	}
	os.WriteFile(clusterPath, data, 0644)
	return nil
}

func initControllerNode(nodeId string, serverConfig *config.ServerConfig) (err error) {
	exist := false // whether the file exist

	// the cluster metata file
	clusterPath := GetClusterFilePath()
	exist, err = util.PathExists(clusterPath)
	if err != nil {
		return err
	}
	if !exist {
		fmt.Println("creating cluster metadata file...")
		var file *os.File
		file, err = os.Create(clusterPath)
		if err != nil {
			return errors.New(cluster_init_failed)
		}
		defer file.Close()

		addrSplit, err := util.SplitAddress(serverConfig.Listen)
		if err != nil {
			panic(err)
		}

		clusterInfo.Ctls[nodeId] = NodeSlim{
			Id:    nodeId,
			Addr:  serverConfig.Listen,
			Addr2: addrSplit[0] + ":" + addrSplit[2],
		}
		var jsonBytes []byte
		jsonBytes, err = json.Marshal(clusterInfo)
		if err != nil {
			fmt.Printf("cluster_init::initControllerNode() %s\n", err.Error())
			panic(err)
		}
		file.WriteString(string(jsonBytes))
	}

	return nil
}
