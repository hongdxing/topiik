/*
* @author: Duan Hongxing
* @date: 23 Aug, 2024
* @desc:
*	Partition info implementation
 */

package cluster

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"topiik/internal/consts"
	"topiik/internal/util"
	"topiik/node"
)

/*
* Client issue command NEW-PARTITION
*
 */
func NewPartition(ptnCount int) (ptnIds []string, err error) {
	if len(partitionInfo.ClusterId) == 0 { // brand new cluster without partition yet
		for i := 0; i < int(ptnCount); i++ {
			var from int
			var to int
			from = i * (consts.SLOTS / ptnCount) // p=2--> i=0: 0, i=1: 512

			if i == (ptnCount - 1) {
				to = consts.SLOTS - 1
			} else {
				to = (i+1)*(consts.SLOTS/ptnCount) - 1 // p=2--> i=0: 511, i=1: 1024
			}
			slot := node.Slot{From: uint16(from), To: uint16(to)}

			ptnId := strings.ToLower(util.RandStringRunes(10))
			ptnIds = append(ptnIds, ptnId)
			partitionInfo.ClusterId = controllerInfo.ClusterId
			partitionInfo.PtnMap[ptnId] = &node.Partition{
				Id:      ptnId,
				NodeSet: make(map[string]*node.NodeSlim),
				Slots:   []node.Slot{slot},
			}
		}
	} else { // having existing partition(s), TODO
		err = errors.New("cannot create new partition")
		return ptnIds, err
	}

	// persist
	err = savePartition()

	return ptnIds, err
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
