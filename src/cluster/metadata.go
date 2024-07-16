/***
** author: duan hongxing
** data: 6 Jul 2024
** desc:
**
**/

package cluster

import (
	"encoding/json"
	"fmt"
	"os"
	"topiik/internal/util"
)

// var metadata = Metadata{}
var nodeInfo *Node
var clusterInfo = &Cluster{Controllers: make(map[string]NodeSlim), Workers: make(map[string]NodeSlim)}
var nodeStatus = &NodeStatus{Role: RAFT_FOLLOWER, Term: 0}

const (
	slash   = string(os.PathSeparator)
	dataDIR = "data"
)

func LoadControllerMetadata(node *Node) (err error) {

	nodeInfo = node

	exist := false // whether the file exist

	// the cluster file
	clusterPath := GetClusterFilePath()
	exist, err = util.PathExists(clusterPath)
	if err != nil {
		panic(err)
	}
	if exist {
		fmt.Println("loading cluster metadata...")
		jsonStr, err := os.ReadFile(clusterPath)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(jsonStr), &clusterInfo)
		fmt.Println(clusterInfo)
		if err != nil {
			panic(err)
		}
	}

	if len(clusterInfo.Controllers) >= 1 {
		go RequestVote()
	}

	fmt.Printf("current node role: %d\n", nodeStatus.Role)

	return nil
}

func UpdateNodeClusterId(clusterId string) (err error) {
	nodeInfo.ClusterId = clusterId
	err = os.Truncate(GetNodeFilePath(), 0)
	if err != nil {
		return err
	}
	buf, err := json.Marshal(nodeInfo)
	if err != nil {
		return err
	}
	err = os.WriteFile(GetNodeFilePath(), buf, 0644)
	if err != nil {
		return err
	}
	return nil
}

func IsNodeController() bool {
	// if current node controllerMap has value
	return len(clusterInfo.Controllers) > 0
}

func GetNodeInfo() Node {
	return *nodeInfo
}

func GetNodeStatus() NodeStatus {
	return *nodeStatus
}

// metadata file path
func GetClusterFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__cluster_metadata"
}
func GetNodeFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__node_metadata"
}
