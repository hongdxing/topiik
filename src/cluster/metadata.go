//author: duan hongxing
//data: 6 Jul 2024

package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"topiik/internal/consts"
	"topiik/internal/util"
	"topiik/node"
)

// var clusterInfo = &Cluster{Ctls: make(map[string]node.NodeSlim), Wkrs: make(map[string]node.NodeSlim)}
var term int

var workerGroupInfo = &WorkerGroupInfo{Groups: make(map[string]*WorkerGroup)}

// var controllerInfo = &NodesInfo{Nodes: make(map[string]node.NodeSlim)}
var partitionInfo = &PartitionInfo{PtnMap: make(map[string]*node.Partition), Slots: make(map[uint16]string, consts.SLOTS)}
var nodeStatus = &NodeStatus{Role: RAFT_FOLLOWER, Term: 0}

const (
	slash   = string(os.PathSeparator)
	dataDIR = "data"
)

// Load controller info on each node, including controller and worker
func LoadWorkerGroupInfo() (err error) {
	l.Info().Msg("Loading controller info begin")
	// whether the file exist
	exist := false

	// Load controller info
	fpath := getWorkerGroupFilePath()
	exist, err = util.PathExists(fpath)
	if err != nil {
		l.Panic().Msg(err.Error())
	}
	if exist {
		jsonStr, err := os.ReadFile(fpath)
		if err != nil {
			l.Panic().Msg(err.Error())
		}
		err = json.Unmarshal([]byte(jsonStr), &workerGroupInfo)
		if err != nil {
			l.Panic().Msg(err.Error())
		}
	}
	l.Info().Msg("Loading controller info end")
	return nil
}

func getWorkGroup(ndId string) WorkerGroup {
	for _, wg := range workerGroupInfo.Groups {
		if _, ok := wg.Nodes[ndId]; ok {
			return *wg
		}
	}
	return WorkerGroup{}
}

func GetNodeStatus() NodeStatus {
	return *nodeStatus
}

func GetTerm() int {
	return term
}

func GetNodeByKeyHash(keyHash uint16) (node.NodeSlim, error) {
	for _, group := range workerGroupInfo.Groups {
		if _, ok := group.Slots[keyHash]; ok {
			return group.Nodes[group.LeaderNodeId], nil
		}
	}
	return node.NodeSlim{}, errors.New(fmt.Sprintf("Cannot find worker for key hash %v", keyHash))
}

func getWorkerGroupFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__metadata_wg__"
}

func getWorkerFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__metadata_worker__"
}

func notifyWorkerGroupChanged() {
	l.Info().Msg("metadata::notifyWorkerGroupChanged begin")
	wrkGrpUpdCh <- struct{}{}
	l.Info().Msg("metadata::notifyWorkerGroupChanged end")
}
