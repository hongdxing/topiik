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
	"io"
	"os"
	"topiik/internal/util"
)

const cluster_init_failed = "cluster init failed"

func ClusterInit(addr string) (err error) {
	fmt.Println("ClusterInit start...")
	// 1. open node file
	nodePath := GetNodeFilePath()
	/*var f *os.File
	f, err = os.OpenFile(nodePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return errors.New(cluster_init_failed)
	}
	defer f.Close()*/

	// 2. read from node file
	var jsonStr string
	var buf []byte
	/*buf := make([]byte, 128)
	var n int
	for {
		n, err = f.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		jsonStr += string(buf[:n])
	}
	if len(jsonStr) == 0 {
		return errors.New(cluster_init_failed)
	}*/

	buf, err = os.ReadFile(nodePath)
	if err != nil{
		return errors.New(cluster_init_failed)
	}
	jsonStr = string(buf)

	// 3. unmarshal
	var node Node
	err = json.Unmarshal([]byte(jsonStr), &node)
	if err != nil {
		return errors.New(cluster_init_failed)
	}
	node.ClusterId = util.RandStringRunes(16)

	// 4. marshal
	buf2, _ := json.Marshal(node)

	// 5. write back to file
	/*_, err = f.Write(buf2)
	if err != nil {
		return errors.New(cluster_init_failed)
	}*/
	os.Truncate(nodePath, 0)
	os.WriteFile(nodePath, buf2, 0644)

	nodeStatus.Role = CCSS_ROLE_CA
	err = initControllerNode()
	if err != nil {
		return err
	}
	go StartServer(addr)
	fmt.Println("ClusterInit end")
	return nil
}

func initControllerNode() (err error) {
	exist := false // whether the file exist

	//var captialMap = make(map[string]Controller)
	var workerMap = make(map[string]Worker)
	var partitionMap = make(map[string]Partition)

	// the controller file
	/*controllerPath := GetControllerFilePath()
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

		controller := ccss.Controller{
			Id:      nodeId,
			Address: serverConfig.Listen,
		}
		captialMap[nodeId] = controller
		var jsonBytes []byte
		jsonBytes, err = json.Marshal(captialMap)
		file.WriteString(string(jsonBytes))
	} else {
		fmt.Println("loading controller metadata...")
		var file *os.File
		file, err = os.Open(controllerPath)
		if err != nil {
			return errors.New(cluster_init_failed)
		}
		defer file.Close()

		captialMap = readMetadata[map[string]Controller](*file)
		fmt.Println(captialMap)
	}*/

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

		workerMap = readMetadata[map[string]Worker](*file)
		fmt.Println(workerMap)
	}

	// the partition file
	patitionPath := GetPartitionFilePath()
	exist, err = util.PathExists(patitionPath)
	if err != nil {
		return err
	}
	if !exist {
		fmt.Println("creating partition file...")
		var file *os.File
		file, err = os.Create(patitionPath)
		if err != nil {
			return errors.New(cluster_init_failed)
		}
		defer file.Close()
	} else {
		fmt.Println("loading partition metadata...")
		var file *os.File
		file, err = os.Open(patitionPath)
		if err != nil {
			return errors.New(cluster_init_failed)
		}
		defer file.Close()

		partitionMap = readMetadata[map[string]Partition](*file)
		fmt.Println(partitionMap)
	}

	return nil
}

func readMetadata[T any](file os.File) (t T) {
	var jsonBytes = make([]byte, 512)
	var jsonStr string
	for {
		n, err := file.Read(jsonBytes)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		jsonStr += string(jsonBytes[:n])
	}
	if len(jsonStr) > 0 {
		err := json.Unmarshal([]byte(jsonStr), &t)
		if err != nil {
			panic(err)
		}
	}
	return t
}
