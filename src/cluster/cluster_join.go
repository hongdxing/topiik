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
	"topiik/internal/consts"
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

		if exist, ok := (controllerMap)[id]; ok {
			//return "", errors.New(node_alread_in_cluster + nodeInfo.ClusterId)
			if exist.Address == addr {
				return nodeInfo.ClusterId, nil
			} else {
				exist.Address = addr // update adddress
			}
		} else {
			controller := Controller{
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
		}

	} else if strings.ToUpper(role) == ROLE_WORKER {

		addrSplit, err := util.SplitAddress(addr)
		if err != nil {
			return "", errors.New(consts.RES_INVALID_ADDR)
		}
		worker := Worker{
			Id:       id,
			Address:  addr,
			Address2: addrSplit[0] + ":" + addrSplit[2],
		}
		if _, ok := workerMap[worker.Id]; ok {
			return "", errors.New(node_alread_in_cluster + nodeInfo.ClusterId)
		}
		workerMap[worker.Id] = worker
		fmt.Println(workerMap)
		// add conttoller id to pending append map
		for _, v := range workerMap {
			if v.Id != nodeInfo.Id {
				workerPendingAppend[v.Id] = v.Id
			}
		}
	} else {
		fmt.Printf("err: %s\n", pieces)
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	fmt.Println(nodeInfo.ClusterId)
	return nodeInfo.ClusterId, nil
}
