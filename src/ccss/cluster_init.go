/***
**
**
**
**
**/

package ccss

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"topiik/internal/config"
	"topiik/internal/util"
)

const cluster_init_failed = "cluster init failed"

func ClusterInit(serverConfig *config.ServerConfig) (err error) {
	fmt.Println("ClusterInit start...")
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
	var node Node
	err = json.Unmarshal([]byte(jsonStr), &node)
	if err != nil {
		return errors.New(cluster_init_failed)
	}
	if len(node.ClusterId) > 0 { // check if current node already in cluster or not
		return errors.New("current node already in a cluster:" + node.ClusterId)
	}
	node.ClusterId = util.RandStringRunes(16)

	// 4. marshal
	buf2, _ := json.Marshal(node)

	// 5. write back to file
	os.Truncate(nodePath, 0)
	os.WriteFile(nodePath, buf2, 0644)

	nodeStatus.Role = RAFT_LEADER
	err = initControllerNode(node.Id, serverConfig)
	if err != nil {
		return err
	}
	go RequestVote()
	fmt.Println("ClusterInit end")
	return nil
}

func initControllerNode(nodeId string, serverConfig *config.ServerConfig) (err error) {
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
	if !exist {
		fmt.Println("creating controller file...")
		var file *os.File
		file, err = os.Create(controllerPath)
		if err != nil {
			return errors.New(cluster_init_failed)
		}
		defer file.Close()

		controller := Controller{
			Id:      nodeId,
			Address: serverConfig.Listen,
		}
		controllerMap[nodeId] = controller
		var jsonBytes []byte
		jsonBytes, err = json.Marshal(controllerMap)
		file.WriteString(string(jsonBytes))
	} else {
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
	if !exist {
		fmt.Println("creating worker file...")
		var file *os.File
		file, err = os.Create(workerPath)
		if err != nil {
			return errors.New(cluster_init_failed)
		}
		defer file.Close()
	} else {
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
	if !exist {
		fmt.Println("creating partition file...")
		var file *os.File
		file, err = os.Create(partitionPath)
		if err != nil {
			return errors.New(cluster_init_failed)
		}
		defer file.Close()
	} else {
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
