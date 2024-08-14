/***
** author: duan hongxing
** data: 4 Jul 2024
** desc:
**
**/

package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

/***
** Obsoleted
** Parameter:
**	- pieces: id[0] host:port[1] ROLE[2]
**
** Syntax: CLUSTER JOIN_ACK id host:port ROLE
**	id: the node id asking to join
**	host:port: the worker listen host:port
**
**/
func clusterJoin(pieces []string) (result string, err error) {
	fmt.Printf("%s\n", pieces)
	if len(pieces) != 3 {
		return "", errors.New(resp.RES_SYNTAX_ERROR)
	}
	id := pieces[0]
	addr := pieces[1]
	role := pieces[2]
	if strings.ToUpper(role) == ROLE_CONTROLLER { // add controller
		var addrSplit []string
		addrSplit, err = util.SplitAddress(addr)
		if err != nil {
			return "", errors.New("join cluster failed")
		}

		if exist, ok := clusterInfo.Ctls[id]; ok {
			if exist.Addr == addr {
				//return nodeInfo.ClusterId, nil
				return "", nil
			} else {
				exist.Addr = addr // update adddress
			}
		} else {
			clusterInfo.Ctls[id] = node.NodeSlim{
				Id:    id,
				Addr:  addr,
				Addr2: addrSplit[0] + ":" + addrSplit[2],
			}
		}

	} else if strings.ToUpper(role) == ROLE_WORKER { // add worker
		var addrSplit []string
		addrSplit, err = util.SplitAddress(addr)
		if err != nil {
			return "", errors.New("join cluster failed")
		}

		if exist, ok := clusterInfo.Wkrs[id]; ok {
			if exist.Addr == addr {
				//return nodeInfo.ClusterId, nil
				return "", nil
			} else {
				exist.Addr = addr // update adddress
			}
		} else {
			clusterInfo.Wkrs[id] = Worker{
				Id:    id,
				Addr:  addr,
				Addr2: addrSplit[0] + ":" + addrSplit[2],
			}
		}
	} else {
		fmt.Printf("err: %s\n", pieces)
		return "", errors.New(resp.RES_SYNTAX_ERROR)
	}

	// increase version
	clusterInfo.Ver += 1

	// persist cluster metadata
	clusterPath := GetClusterFilePath()
	buf, err := json.Marshal(clusterInfo)
	if err != nil {
		return "", errors.New("update cluster failed")
	}
	err = os.Truncate(clusterPath, 0) // TODO: myabe need backup first
	if err != nil {
		return "", errors.New("update cluster failed")
	}
	err = os.WriteFile(clusterPath, buf, 0664) // save back controller file
	if err != nil {
		return "", errors.New("update cluster failed")
	}

	// cluster meta changed, pending to sync to follower(s)
	/*
		for _, v := range clusterInfo.Ctls {
			if v.Id != nodeInfo.Id {
				clusterMetadataPendingAppend[v.Id] = v.Id
			}
		}
	*/
	cluUpdCh <- struct{}{}
	//fmt.Println(nodeInfo.ClusterId)
	//return nodeInfo.ClusterId, nil
	return "", nil
}
