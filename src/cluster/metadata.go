/*
* author: duan hongxing
* data: 6 Jul 2024
* desc:
*
 */

package cluster

import (
	"encoding/json"
	"fmt"
	"os"
	"topiik/internal/util"
	"topiik/node"
)

// var clusterInfo = &Cluster{Ctls: make(map[string]node.NodeSlim), Wkrs: make(map[string]node.NodeSlim)}
var term int
var controllerInfo = &NodesInfo{Nodes: make(map[string]node.NodeSlim)}
var workerInfo = &NodesInfo{Nodes: make(map[string]node.NodeSlim)}
var partitionInfo = &PartitionInfo{PtnMap: make(map[string]*node.Partition)}
var nodeStatus = &NodeStatus{Role: RAFT_FOLLOWER, Term: 0}

const (
	slash   = string(os.PathSeparator)
	dataDIR = "data"
)

func LoadControllerMetadata() (err error) {
	exist := false // whether the file exist

	// the cluster file
	fpath := getControllerFilePath()
	exist, err = util.PathExists(fpath)
	if err != nil {
		l.Panic().Msg(err.Error())
	}
	if exist {
		l.Info().Msg("Loading cluster metadata...")
		jsonStr, err := os.ReadFile(fpath)
		if err != nil {
			l.Panic().Msg(err.Error())
		}
		err = json.Unmarshal([]byte(jsonStr), &controllerInfo)
		fmt.Println(controllerInfo)
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
	l.Info().Msgf("Current node role: %s", node.GetNodeInfo().Role)
	if node.GetNodeInfo().Role == node.ROLE_CONTROLLER {
		go RequestVote()
	}

	return nil
}

func SetControllerInfo(controllers *NodesInfo) {
	controllerInfo = controllers
	saveControllerInfo()
}

func SetWorkerInfo(workers *NodesInfo) {
	workerInfo = workers
	saveWorkerInfo()
}

/*
	func GetClusterInfo() Cluster {
		return *clusterInfo
	}
*/
func GetWorkerInfo() NodesInfo {
	return *workerInfo
}

func GetWorkerLeaders() (workers []node.NodeSlim) {
	for _, ptn := range partitionInfo.PtnMap {
		worker := workerInfo.Nodes[ptn.LeaderNodeId]
		workers = append(workers, worker)
	}
	return workers
}

func GetNodeStatus() NodeStatus {
	return *nodeStatus
}

func GetTerm() int {
	return term
}

/* metadata file path */
/*
func GetClusterFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__metadata_cluster__"
}*/

func getControllerFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__metadata_controller__"
}

func getWorkerFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__metadata_worker__"
}

func GetPatitionFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__metadata_partition__"
}

/* save meatadata */
func saveControllerInfo() (err error) {
	data, err := json.Marshal(controllerInfo)
	if err != nil {
		l.Err(err).Msgf("cluster::saveControllerInfo %s", err.Error())
		return err
	}

	fpath := getControllerFilePath()
	exist, _ := util.PathExists(fpath)
	if exist {
		err = os.Truncate(fpath, 0) // TODO: backup first
		if err != nil {
			l.Err(err)
			return err
		}
	}

	err = os.WriteFile(fpath, data, 0644)
	if err != nil {
		l.Err(err)
		return err
	}
	return nil
}

func saveWorkerInfo() (err error) {
	data, err := json.Marshal(workerInfo)
	if err != nil {
		l.Err(err).Msgf("cluster::saveWorkerInfo %s", err.Error())
		return err
	}

	fpath := getWorkerFilePath()
	exist, _ := util.PathExists(fpath)
	if exist {
		err = os.Truncate(fpath, 0) // TODO: backup first
		if err != nil {
			l.Err(err)
			return err
		}
	}

	err = os.WriteFile(fpath, data, 0644)
	if err != nil {
		l.Err(err)
		return err
	}
	return nil
}

func notifyMetadataChanged() {
	l.Info().Msg("metadata::notifyMetadataChanged begin")
	cluUpdCh <- struct{}{}
	l.Info().Msg("metadata::notifyMetadataChanged end")
}

func notifyControllerChanged() {
	l.Info().Msg("metadata::notifyControllerChanged begin")
	ctlUpdCh <- struct{}{}
	l.Info().Msg("metadata::notifyControllerChanged end")
}

func notifyWorkerChanged() {
	l.Info().Msg("metadata::notifyWorkerChanged begin")
	wrkUpdCh <- struct{}{}
	l.Info().Msg("metadata::notifyWorkerChanged end")
}

func notifyPtnChanged() {
	l.Info().Msg("metadata::notifyPtnChange begin")
	ptnUpdCh <- struct{}{}
	l.Info().Msg("metadata::notifyPtnChange end")
}
