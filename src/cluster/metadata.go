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


var clusterInfo = &Cluster{Ctls: make(map[string]NodeSlim), Wkrs: make(map[string]Worker)}
var partitionInfo = &PartitionInfo{PtnMap: make(map[string]Partition)}
var nodeStatus = &NodeStatus{Role: RAFT_FOLLOWER, Term: 0}


const (
	slash   = string(os.PathSeparator)
	dataDIR = "data"
)

func LoadControllerMetadata() (err error) {
	exist := false // whether the file exist

	// the cluster file
	clusterPath := GetClusterFilePath()
	exist, err = util.PathExists(clusterPath)
	if err != nil {
		l.Panic().Msg(err.Error())
	}
	if exist {
		l.Info().Msg("Loading cluster metadata...")
		jsonStr, err := os.ReadFile(clusterPath)
		if err != nil {
			l.Panic().Msg(err.Error())
		}
		err = json.Unmarshal([]byte(jsonStr), &clusterInfo)
		fmt.Println(clusterInfo)
		if err != nil {
			l.Panic().Msg(err.Error())
			//panic(err)
		}
	}

	// the slots file

	filePath := GetPatitionFilePath()
	exist, err = util.PathExists(filePath)
	if err != nil {
		l.Panic().Msg(err.Error())
	}
	if exist {
		data, err := util.ReadBinaryFile(filePath)
		if err != nil {
			l.Panic().Msg(err.Error())
		}

		if err != nil {
			l.Panic().Msg(err.Error())
		}
		err = json.Unmarshal(data, &partitionInfo)
		if err != nil {
			l.Panic().Msg(err.Error())
		}
	}

	//
	l.Info().Msgf("Current node role: %d", nodeStatus.Role)
	if len(clusterInfo.Ctls) >= 1 {
		go RequestVote()
	}

	return nil
}

func UpdatePendingAppend() {
	l.Info().Msg("metadata::UpdatePendingAppend begin")
	cluUpdCh <- struct{}{}
	l.Info().Msg("metadata::UpdatePendingAppend end")
}

func IsNodeController() bool {
	// if current node controllerMap has value
	return len(clusterInfo.Ctls) > 0
}


func GetClusterInfo() Cluster {
	return *clusterInfo
}

func GetWorkerLeaders() (workers []Worker) {
	for _, ptn := range partitionInfo.PtnMap {
		worker := clusterInfo.Wkrs[ptn.LeaderNodeId]
		workers = append(workers, worker)
	}
	return workers
}

func GetNodeStatus() NodeStatus {
	return *nodeStatus
}

// metadata file path
func GetClusterFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__metadata_cluster__"
}

func GetPatitionFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__metadata_partition__"
}
