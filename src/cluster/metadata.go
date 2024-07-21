/***
** author: duan hongxing
** data: 6 Jul 2024
** desc:
**
**/

package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"topiik/internal/util"
)

// var metadata = Metadata{}
var nodeInfo *Node
var clusterInfo = &Cluster{Ctls: make(map[string]NodeSlim), Wkrs: make(map[string]Worker)}
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
		tLog.Panic().Msg(err.Error())
	}
	if exist {
		tLog.Info().Msg("Loading cluster metadata...")
		jsonStr, err := os.ReadFile(clusterPath)
		if err != nil {
			tLog.Panic().Msg(err.Error())
		}
		err = json.Unmarshal([]byte(jsonStr), &clusterInfo)
		fmt.Println(clusterInfo)
		if err != nil {
			tLog.Panic().Msg(err.Error())
			//panic(err)
		}
	}

	tLog.Info().Msgf("Current node role: %d", nodeStatus.Role)
	if len(clusterInfo.Ctls) >= 1 {
		go RequestVote()
	}

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

func AddNode(nodeId string, addr string, addr2 string, role string) (err error) {
	if strings.ToUpper(role) == ROLE_CONTROLLER {
		clusterInfo.Ctls[nodeId] = NodeSlim{Id: nodeId, Addr: addr, Addr2: addr2}
	} else {
		clusterInfo.Wkrs[nodeId] = Worker{Id: nodeId, Addr: addr, Addr2: addr2}
	}
	clusterPath := GetClusterFilePath()
	buf, err := json.Marshal(clusterInfo)
	if err != nil {
		return errors.New("update cluster failed")
	}
	err = os.Truncate(clusterPath, 0) // TODO: myabe need backup first
	if err != nil {
		return errors.New("update cluster failed")
	}
	err = os.WriteFile(clusterPath, buf, 0664) // save back controller file
	if err != nil {
		return errors.New("update cluster failed")
	}

	return nil
}

func UpdatePendingAppend() {
	for _, v := range clusterInfo.Ctls {
		if v.Id != nodeInfo.Id {
			clusterMetadataPendingAppend[v.Id] = v.Id
		}
	}
}

func IsNodeController() bool {
	// if current node controllerMap has value
	return len(clusterInfo.Ctls) > 0
}

func GetNodeInfo() Node {
	return *nodeInfo
}

func GetNodeStatus() NodeStatus {
	return *nodeStatus
}

// metadata file path
func GetClusterFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__metadata_cluster__"
}
func GetNodeFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__metadata_node__"
}
