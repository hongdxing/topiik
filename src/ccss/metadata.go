/***
** author: duan hongxing
** data: 6 Jul 2024
** desc:
**
**/

package ccss

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"topiik/internal/util"
)

var meatadata = Metadata{}
var nodeStatus = &NodeStatus{Role: RAFT_FOLLOWER, Term: 0}
var controllerMap = make(map[string]Controller)
var workerMap = make(map[string]Worker)
var partitionMap = make(map[string]Partition)

func InitMetadata(controllers map[string]Controller, solars map[string]Worker, partitions map[string]Partition) {
	controllerMap = controllers
	fmt.Println("controllers")
	fmt.Println(controllerMap)
	workerMap = solars
	fmt.Println("workers")
	fmt.Println(workerMap)
	partitionMap = partitions
	fmt.Println("partitions")
	fmt.Println(partitionMap)
}

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

	meatadata.Node = *node

	exist := false // whether the file exist

	//var captialMap = make(map[string]Controller)
	var workerMap = make(map[string]Worker)
	var partitionMap = make(map[string]Partition)

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

		controllerMap = readMetadata[map[string]Controller](controllerPath)
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

		workerMap = readMetadata[map[string]Worker](workerPath)
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

		partitionMap = readMetadata[map[string]Partition](partitionPath)
		fmt.Println(partitionMap)
	}

	return nil
}

func readMetadata[T any](metadataPath string) (t T) {
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
	return t
}

func UpdateNodeClusterId(clusterId string) (err error) {
	fmt.Println(meatadata)
	node := meatadata.Node
	node.ClusterId = clusterId
	err = os.Truncate(GetNodeFilePath(), 0)
	if err != nil {
		return err
	}
	buf, err := json.Marshal(node)
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
	if len(controllerMap) > 0 {
		return true
	}
	return false
}

func GetNodeMetadata() Node {
	return meatadata.Node
}

// metadata file path
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
