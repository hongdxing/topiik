/*
* author: duan hongxing
* date: 6 July, 2024
* desc:
*
 */

package cluster

import (
	"encoding/json"
	"errors"
	"os"
	"topiik/internal/config"
	"topiik/internal/util"
	"topiik/node"
)

const clusterInitFailed = "cluster init failed"

/*
* Execute command INIT-CLUSTER
*
 */
func InitCluster(serverConfig *config.ServerConfig) error {
	l.Info().Msg("cluster::ClusterInit start")

	// 0. init cluster
	err := doInit(serverConfig)
	if err != nil {
		return err
	}
	node.InitCluster(GetClusterInfo().Id)
	nodeStatus.Role = RAFT_LEADER
	/*err = initControllerNode(nodeId, serverConfig)
	if err != nil {
		return err
	}*/
	// after init, the node default is LEADER, and start to AppendEntries()
	go AppendEntries()
	l.Info().Msg("cluster::ClusterInit end")
	return nil
}

func doInit(serverConfig *config.ServerConfig) error {
	if len(clusterInfo.Id) > 0 {
		return errors.New("current node already in cluster: " + clusterInfo.Id)
	}
	// set clusterInfo
	nodeId := node.GetNodeInfo().Id
	clusterInfo.Id = util.RandStringRunes(10)

	hostPort, err := util.SplitAddress(serverConfig.Listen)
	if err != nil {
		l.Panic().Msg(err.Error())
	}
	clusterInfo.Ctls[nodeId] = node.NodeSlim{
		Id:    nodeId,
		Addr:  serverConfig.Listen,
		Addr2: hostPort[0] + ":" + hostPort[2],
	}

	// save cluster metadata
	fpath := GetClusterFilePath()
	_ = os.Remove(fpath) // just incase
	data, err := json.Marshal(clusterInfo)
	if err != nil {
		return errors.New(clusterInitFailed)
	}
	err = os.WriteFile(fpath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

/*
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
}*/
