// author: Duan Hongxing
// date: 23 Aug, 2024
// desc:
//	Partition info implementation

package cluster

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"topiik/internal/consts"
	"topiik/internal/util"
	"topiik/node"
)

// New partition
func NewPartition(workers map[string]string) (ptnId string, err error) {

	//for i := 0; i < int(ptnCount); i++ {
	//	ptnId := strings.ToLower(util.RandStringRunes(10))
	//	partitionInfo.PtnMap[ptnId] = &node.Partition{
	//		Id:      ptnId,
	//		NodeSet: make(map[string]*node.NodeSlim),
	//		// The SlotFrom and SlotTo leave to not set, Till RESHARD executed
	//	}
	//}

	addPartition(workers)

	return ptnId, nil
}

// add partition
// controllers: map of node id and addr
func addPartition(controllers map[string]string) {
	ptnId := strings.ToLower(util.RandStringRunes(10))
	newPartition := &node.Partition{
		Id:      ptnId,
		NodeSet: make(map[string]*node.NodeSlim),
		// The SlotFrom and SlotTo leave to not set, Till RESHARD executed
	}
	i := 0
	for ndId := range controllers {
		// set first node to leader
		if i == 0 {
			newPartition.LeaderNodeId = ndId
		}
		newPartition.NodeSet[ndId] = &node.NodeSlim{Id: ndId}
		i++
	}
	partitionInfo.PtnMap[ptnId] = newPartition
}

// Remove partition
func RemovePartition(ptnId string) error {
	return nil
}

// Reshard partition
func ReShard(isCreate bool) (err error) {
	// brand new cluster without partition yet
	if isCreate {
		var ptnCount = len(workerGroupInfo.Groups)
		var i int = 0
		for _, group := range workerGroupInfo.Groups {
			group.Slots = map[uint16]bool{}
			var from int
			var to int
			from = i * (consts.SLOTS / ptnCount) // p=2--> i=0: 0, i=1: 512
			if i == (ptnCount - 1) {
				to = consts.SLOTS - 1
			} else {
				to = (i+1)*(consts.SLOTS/ptnCount) - 1 // p=2--> i=0: 511, i=1: 1024
			}
			for i := from; i < to; i++ {
				group.Slots[uint16(i)] = true
			}
			i++
		}
		fmt.Println(workerGroupInfo)
	} else {
		//
	}

	saveWorkerGroups()

	return err
}

// save worker group
func saveWorkerGroups() (err error) {
	data, err := json.Marshal(workerGroupInfo)
	if err != nil {
		l.Err(err).Msgf("cluster::saveControllerInfo %s", err.Error())
		return err
	}

	fpath := getWorkerGroupFilePath()
	exist, _ := util.PathExists(fpath)
	if exist { // rename to old for backup
		err = os.Rename(fpath, fpath+"old")
		if err != nil {
			l.Err(err).Msgf("saveWorkerGroups: %s", err.Error())
			return err
		}
	}

	err = util.WriteBinaryFile(fpath, data)
	if err != nil {
		l.Err(err).Msgf("saveWorkerGroups: %s", err.Error())
		return err
	}
	return nil
}
