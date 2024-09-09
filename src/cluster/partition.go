// author: Duan Hongxing
// date: 23 Aug, 2024
// desc:
//	Partition info implementation

package cluster

import (
	"encoding/json"
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
// workers: map of node id and addr
func addPartition(workers map[string]string) {
	ptnId := strings.ToLower(util.RandStringRunes(10))
	newPartition := &node.Partition{
		Id:      ptnId,
		NodeSet: make(map[string]*node.NodeSlim),
		// The SlotFrom and SlotTo leave to not set, Till RESHARD executed
	}
	for ndId := range workers {
		newPartition.NodeSet[ndId] = &node.NodeSlim{Id: ndId}
	}
	partitionInfo.PtnMap[ptnId] = newPartition
}

// Remove partition
func RemovePartition(ptnId string) error {
	return nil
}

// Reshard partition
// The new Partition always take slots from the last slot, i.e. the one end with 1023
// And then slide from the first to next make slots even
func ReShard() (err error) {
	// brand new cluster without partition yet
	if len(partitionInfo.ClusterId) == 0 {
		partitionInfo.ClusterId = controllerInfo.ClusterId
		var ptnCount = len(partitionInfo.PtnMap)
		var i int = 0
		for _, ptn := range partitionInfo.PtnMap {
			var from int
			var to int
			from = i * (consts.SLOTS / ptnCount) // p=2--> i=0: 0, i=1: 512
			if i == (ptnCount - 1) {
				to = consts.SLOTS - 1
			} else {
				to = (i+1)*(consts.SLOTS/ptnCount) - 1 // p=2--> i=0: 511, i=1: 1024
			}
			ptn.SlotFrom = uint16(from)
			ptn.SlotTo = uint16(to)
			i++
		}
	} else {
		//
	}

	// Notify
	notifyPtnChanged()

	// Persist
	err = savePartition()
	return err
}

// Add node to NodeSet of partition
func addNode2Partition(ptnId string, ndId string) {
	if ptn, ok := partitionInfo.PtnMap[ptnId]; ok {
		if len(ptn.NodeSet) == 0 {
			/* if ndId is the first node, set to Leader */
			ptn.LeaderNodeId = ndId
			ptn.NodeSet[ndId] = &node.NodeSlim{
				Id: ndId,
			}
		} else if _, ok := ptn.NodeSet[ndId]; !ok {
			ptn.NodeSet[ndId] = &node.NodeSlim{
				Id: ndId,
			}
		}
	}
	savePartition()
}

func GetPartitionInfo() PartitionInfo {
	return *partitionInfo
}

// Save partition info to disk
func savePartition() (err error) {
	fpath := GetPatitionFilePath()
	exist, err := util.PathExists(fpath)
	if err != nil {
		l.Err(err).Msgf("savePartition: %s", err.Error())
		return err
	}
	if exist { // rename to old for backup
		err = os.Rename(fpath, fpath+"old")
		if err != nil {
			l.Err(err).Msgf("savePartition: %s", err.Error())
			return err
		}
	}
	data, err := json.Marshal(partitionInfo)
	if err != nil {
		l.Err(err).Msgf("savePartition: %s", err.Error())
		return err
	}

	err = util.WriteBinaryFile(fpath, data)
	if err != nil {
		l.Err(err).Msgf("savePartition: %s", err.Error())
		return err
	}
	return nil
}
