/*
* author: duan hongxing
* data: 3 Aug 2024
* desc:
*
 */

package cluster

import (
	"errors"
	"topiik/resp"
)

func getPartitionLeader(pieces []string) (res string, err error) {
	if len(pieces) != 1 {
		return "", errors.New(resp.RES_SYNTAX_ERROR)
	}
	nodeId := pieces[0]
	var ptnLeaderId string
	// get partition leader id
	for _, ptn := range partitionInfo.PtnMap {
		if _, ok := ptn.NodeSet[nodeId]; ok {
			ptnLeaderId = ptn.LeaderNodeId
			break
		}
	}
	if len(ptnLeaderId) == 0 {
		return "", errors.New(resp.RES_NIL)
	}
	// get worker use the partition leader id
	if worker, ok := clusterInfo.Wkrs[ptnLeaderId]; ok {
		// return partition leader node id + addr2
		// when the follow get this string, can use fix lenght worker node id to split it
		return worker.Id + worker.Addr2, nil
	}
	return "", errors.New(resp.RES_NIL)
}
