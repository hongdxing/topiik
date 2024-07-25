/*
* author: duan hongxing
* data: 25 Jul 2024
* desc:
*
 */

package cluster

import (
	"errors"
	"topiik/internal/util"
	"topiik/resp"
)

func Scale(p int, r int) (result string, err error) {
	// if no enough worker nodes
	if len(clusterInfo.Wkrs) < p*r || p > SLOTS {
		return "", errors.New(resp.RES_NO_ENOUGH_WORKER)
	}
	if len(partitionInfo) == 0 { // new cluster
		keys := make([]string, 0, len(clusterInfo.Wkrs))
		for k := range clusterInfo.Wkrs {
			keys = append(keys, k)
		}

		for i := 0; i < int(p); i++ {
			works := keys[i*r : (i+1)*r] // 2*2--> i==0: [0:2], i==1: [2:4]
			pId := util.RandStringRunes(16)
			partition := Partition{
				Id:           pId,
				LeaderNodeId: works[0],
			}

			var from int
			var to int
			from = i * (SLOTS / p) // p=2--> i=0: 0, i=1: 512

			if i == (p - 1) {
				to = SLOTS - 1
			} else {
				to = (i+1)*(SLOTS/p) - 1 // p=2--> i=0: 511, i=1: 1024
			}
			slot := Slot{From: uint16(from), To: uint16(to)}
			partition.Slots = []Slot{slot}
			partitionInfo[pId] = partition
		}
	} else if p > len(partitionInfo) { // scale out
		//
	} else { // scale in
		//
	}
	return resp.RES_OK, nil
}
