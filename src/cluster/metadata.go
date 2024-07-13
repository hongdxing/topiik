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
	"topiik/internal/util"
)

// var metadata = Metadata{}
var nodeInfo *Node
var clusterInfo = &Cluster{Controllers: make(map[string]NodeSlim), Workers: make(map[string]NodeSlim)}
var controllerMap = make(map[string]NodeSlim)
var workerMap = make(map[string]NodeSlim)
var partitionMap = make(map[string]Partition)
var nodeStatus = &NodeStatus{Role: RAFT_FOLLOWER, Term: 0}

func Map2Array[T any](theMap map[string]T) (arr []T) {
	for _, v := range theMap {
		arr = append(arr, v)
	}
	return arr
}

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

	// the controller file
	controllerPath := GetControllerFilePath()
	exist, err = util.PathExists(controllerPath)
	if err != nil {
		return err
	}
	if exist {
		fmt.Println("loading controller metadata...")
		var file *os.File
		file, err = os.Open(controllerPath)
		if err != nil {
			return errors.New(cluster_init_failed)
		}
		defer file.Close()

		//controllerMap = readMetadata[map[string]Controller](controllerPath)
		readMetadata[map[string]NodeSlim](controllerPath, &controllerMap)
		fmt.Println(controllerMap)
	}

	// the worker file
	workerPath := GetWorkerFilePath()
	exist, err = util.PathExists(workerPath)
	if err != nil {
		return err
	}
	if exist {
		fmt.Println("loading worker metadata...")
		var file *os.File
		file, err = os.Open(workerPath)
		if err != nil {
			return errors.New(cluster_init_failed)
		}
		defer file.Close()

		//workerMap = readMetadata[map[string]Worker](workerPath)
		readMetadata[map[string]NodeSlim](workerPath, &workerMap)
		fmt.Println(workerMap)
	}

	// the partition file
	partitionPath := GetPartitionFilePath()
	exist, err = util.PathExists(partitionPath)
	if err != nil {
		return err
	}
	if exist {
		fmt.Println("loading partition metadata...")
		var file *os.File
		file, err = os.Open(partitionPath)
		if err != nil {
			return errors.New(cluster_init_failed)
		}
		defer file.Close()

		//partitionMap = readMetadata[map[string]Partition](partitionPath)
		readMetadata[map[string]Partition](partitionPath, &partitionMap)
		fmt.Println(partitionMap)
	}

	// if current node is Controller Node(stop and restarted), start to RequestVote()
	if len(controllerMap) >= 1 {
		go RequestVote()
	}

	fmt.Printf("current node role: %d\n", nodeStatus.Role)

	return nil
}

func readMetadata[T any](metadataPath string, t *T) {
	var jsonStr string
	buf, err := os.ReadFile(metadataPath)
	if err != nil {
		panic(err)
	}

	jsonStr = string(buf)
	if len(jsonStr) > 0 {
		err := json.Unmarshal([]byte(jsonStr), &t)
		if err != nil {
			panic(err)
		}
	}
	//return t
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
	return len(controllerMap) > 0
}

func GetNodeMetadata() Node {
	return *nodeInfo
}

// metadata file path
func GetClusterFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__cluster_metadata"
}
func GetNodeFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "metadata_node"
}

func GetControllerFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "metadata_controller"
}

func GetWorkerFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "metadata_worker"
}

func GetPartitionFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "metadata_partition"
}
