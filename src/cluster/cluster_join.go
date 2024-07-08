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
**	- pieces: JOIN_ACK[0] id[1] host:port[2] ROLE[3]
**
** Syntax: CLUSTER JOIN_ACK id host:port ROLE
**	id: the node id asking to join
**	host:port: the worker listen host:port
**
**/
func clusterJoin(pieces []string) (result string, err error) {
	fmt.Printf("%s\n", pieces)
	if len(pieces) != 4 {
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	if strings.ToUpper(pieces[3]) == ROLE_CONTROLLER {
		var addrSplit []string
		addrSplit, err = util.SplitAddress(pieces[2])
		if err != nil {
			return "", errors.New("join cluster failed")
		}

		if exist, ok := (controllerMap)[pieces[1]]; ok {
			//return "", errors.New(node_alread_in_cluster + nodeInfo.ClusterId)
			if exist.Address == pieces[2] {
				return nodeInfo.ClusterId, nil
			} else {
				exist.Address = pieces[2] // update adddress
			}
		} else {
			controller := Controller{
				Id:       pieces[1],
				Address:  pieces[2],
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
	} else if strings.ToUpper(pieces[3]) == ROLE_WORKER {

		addrSplit, err := util.SplitAddress(pieces[2])
		if err != nil {
			return "", errors.New(consts.RES_INVALID_ADDR)
		}
		worker := Worker{
			Id:       pieces[1],
			Address:  pieces[2],
			Address2: addrSplit[0] + ":" + addrSplit[2],
		}
		if _, ok := workerMap[worker.Id]; ok {
			return "", errors.New(node_alread_in_cluster + nodeInfo.ClusterId)
		}
		workerMap[worker.Id] = worker
		fmt.Println(workerMap)
	} else {
		fmt.Printf("err: %s\n", pieces)
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	fmt.Println(nodeInfo.ClusterId)
	return nodeInfo.ClusterId, nil
}
