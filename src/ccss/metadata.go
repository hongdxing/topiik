/***
** author: duan hongxing
** data: 6 Jul 2024
** desc:
**
**/

package ccss

import (
	"fmt"
	"os"
	"topiik/internal/util"
)

var nodeStatus = &NodeStatus{Role: CCSS_ROLE_CO, Term: 0}
var captialMap = make(map[string]Controller)
var workerMap = make(map[string]Worker)
var partitionMap = make(map[string]Partition)

func InitMetadata(controllers map[string]Controller, solars map[string]Worker, partitions map[string]Partition) {
	captialMap = controllers
	fmt.Println("controllers")
	fmt.Println(captialMap)
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
