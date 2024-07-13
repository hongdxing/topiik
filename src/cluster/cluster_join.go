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
)

const node_alread_in_cluster = "node already in cluster:"

/***
**
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
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	id := pieces[0]
	addr := pieces[1]
	role := pieces[2]
	if strings.ToUpper(role) == ROLE_CONTROLLER {
		var addrSplit []string
		addrSplit, err = util.SplitAddress(addr)
		if err != nil {
			return "", errors.New("join cluster failed")
		}

		if exist, ok := clusterInfo.Controllers[id]; ok {
			if exist.Address == addr {
				return nodeInfo.ClusterId, nil
			} else {
				exist.Address = addr // update adddress
			}
		} else {
			clusterInfo.Controllers[id] = NodeSlim{
				Id:       id,
				Address:  addr,
				Address2: addrSplit[0] + ":" + addrSplit[2],
			}
		}

		/*if exist, ok := (controllerMap)[id]; ok {
			//return "", errors.New(node_alread_in_cluster + nodeInfo.ClusterId)
			if exist.Address == addr {
				return nodeInfo.ClusterId, nil
			} else {
				exist.Address = addr // update adddress
			}
		} else {
			controller := NodeSlim{
				Id:       id,
				Address:  addr,
				Address2: addrSplit[0] + ":" + addrSplit[2],
			}
			controllerMap[controller.Id] = controller
		}
		fmt.Println(controllerMap)
		controllerPath := GetControllerFilePath()
		buf, err := json.Marshal(controllerMap)
		if err != nil {
			return "", errors.New("save controller failed")
		}
		err = os.Truncate(controllerPath, 0) // TODO: myabe need backup first
		if err != nil {
			return "", errors.New("save controller failed")
		}
		err = os.WriteFile(controllerPath, buf, 0664) // save back controller file
		if err != nil {
			return "", errors.New("save controller failed")
		}
		// add conttoller id to pending append map
		for _, v := range controllerMap {
			if v.Id != nodeInfo.Id {
				controllerPendingAppend[v.Id] = v.Id
			}
		}*/

	} else if strings.ToUpper(role) == ROLE_WORKER {
		var addrSplit []string
		addrSplit, err = util.SplitAddress(addr)
		if err != nil {
			return "", errors.New("join cluster failed")
		}

		if exist, ok := clusterInfo.Workers[id]; ok {
			if exist.Address == addr {
				return nodeInfo.ClusterId, nil
			} else {
				exist.Address = addr // update adddress
			}
		} else {
			clusterInfo.Workers[id] = NodeSlim{
				Id:       id,
				Address:  addr,
				Address2: addrSplit[0] + ":" + addrSplit[2],
			}
		}

		/*
			addrSplit, err := util.SplitAddress(addr)
			if err != nil {
				return "", errors.New(consts.RES_INVALID_ADDR)
			}

			if exist, ok := workerMap[id]; ok {
				if exist.Address == addr {
					return nodeInfo.ClusterId, nil
				} else {
					exist.Address = addr // update adddress
				}
				return "", errors.New(node_alread_in_cluster + nodeInfo.ClusterId)
			} else {
				worker := NodeSlim{
					Id:       id,
					Address:  addr,
					Address2: addrSplit[0] + ":" + addrSplit[2],
				}
				workerMap[worker.Id] = worker
			}
			fmt.Println(workerMap)
			workerPath := GetWorkerFilePath()
			buf, err := json.Marshal(workerMap)
			if err != nil {
				return "", errors.New("save worker failed")
			}
			err = os.Truncate(workerPath, 0) // TODO: myabe need backup first
			if err != nil {
				return "", errors.New("save worker failed")
			}
			err = os.WriteFile(workerPath, buf, 0664) // save back worker file
			if err != nil {
				return "", errors.New("save worker failed")
			}
			// add woker id to pending append map
			for _, v := range workerMap {
				if v.Id != nodeInfo.Id {
					workerPendingAppend[v.Id] = v.Id
				}
			}*/
	} else {
		fmt.Printf("err: %s\n", pieces)
		return "", errors.New(RES_SYNTAX_ERROR)
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
	for _, v := range clusterInfo.Controllers {
		if v.Id != nodeInfo.Id {
			clusterMetadataPendingAppend[v.Id] = v.Id
		}
	}
	fmt.Println(nodeInfo.ClusterId)
	return nodeInfo.ClusterId, nil
}
