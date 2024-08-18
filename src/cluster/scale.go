/*
* author: duan hongxing
* data: 25 Jul 2024
* desc:
*
 */

package cluster

import (
	"encoding/json"
	"errors"
	"os"
	"topiik/internal/consts"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

func Scale(p int, r int) (result string, err error) {
	// if no enough worker nodes
	if len(clusterInfo.Wkrs) < p*r || p > consts.SLOTS {
		return "", errors.New(resp.RES_NO_ENOUGH_WORKER)
	}
	if len(partitionInfo.PtnMap) == 0 { // new cluster
		keys := make([]string, 0, len(clusterInfo.Wkrs))
		for k := range clusterInfo.Wkrs {
			keys = append(keys, k)
		}

		for i := 0; i < int(p); i++ {
			works := keys[i*r : (i+1)*r] // 2*2--> i==0: [0:2], i==1: [2:4]
			pId := util.RandStringRunes(consts.PTN_ID_LEN)
			partition := node.Partition{
				Id:           pId,
				LeaderNodeId: works[0],
			}
			nodeSet := make(map[string]*node.NodeSlim)
			for _, worker := range keys {
				nodeSet[worker] = &node.NodeSlim{} //
			}
			partition.NodeSet = nodeSet

			var from int
			var to int
			from = i * (consts.SLOTS / p) // p=2--> i=0: 0, i=1: 512

			if i == (p - 1) {
				to = consts.SLOTS - 1
			} else {
				to = (i+1)*(consts.SLOTS/p) - 1 // p=2--> i=0: 511, i=1: 1024
			}
			slot := node.Slot{From: uint16(from), To: uint16(to)}
			partition.Slots = []node.Slot{slot}
			partitionInfo.Ptns = uint16(p)
			partitionInfo.Rpls = uint16(r)
			partitionInfo.PtnMap[pId] = &partition
		}
	} else if p > len(partitionInfo.PtnMap) { // scale out
		//
	} else { // scale in
		//
	}
	// persist
	filePath := GetPatitionFilePath()
	exist, err := util.PathExists(filePath)
	if err != nil {
		l.Err(err).Msgf("scale: %s", err.Error())
		return "", err
	}
	if exist { // rename to old for backup
		os.Rename(filePath, filePath+"old")
	}
	data, err := json.Marshal(partitionInfo)
	if err != nil {
		l.Err(err).Msgf("scale: %s", err.Error())
		return "", err
	}

	err = util.WriteBinaryFile(filePath, data)
	if err != nil {
		l.Err(err).Msgf("scale: %s", err.Error())
		return "", err
	}

	// update cluster info
	//clusterInfo.Ptns = uint16(p)
	//clusterInfo.Rpls = uint16(r)
	//UpdatePendingAppend()  // sync cluster info to followers
	ptnUpdCh <- struct{}{} // sync partition to followers

	return resp.RES_OK, nil
}
