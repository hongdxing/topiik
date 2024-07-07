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

/***
**
** Parameter:
**	- pieces: JOIN_ACK id host:port ROLE
**
** Syntax: CLUSTER JOIN_ACK id host:port ROLE
**	id: the node id asking to join
**	host:port: the worker listen host:port
**
**/
func clusterJoin(pieces []string) (result string, err error) {
	if len(pieces) != 4 {
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	fmt.Printf(" %s\n", pieces)
	if strings.ToUpper(pieces[3]) == ROLE_CONTROLLER {
		var addrSplit []string
		addrSplit, err = util.SplitAddress(pieces[2])
		if err != nil {
			return "", errors.New("join cluster failed")
		}
		controller := Controller{
			Id:       pieces[1],
			Address:  pieces[2],
			Address2: addrSplit[0] + ":" + addrSplit[2],
		}
		if _, ok := (controllerMap)[controller.Id]; ok {
			return "", errors.New("WORKER_ALREADY_IN_CLUSTER:" + controller.Id)
		}
		(controllerMap)[controller.Id] = controller
		fmt.Println(controllerMap)
		controllerPath := GetControllerFilePath()
		err = os.Truncate(controllerPath, 0) // TODO: myabe need backup first
		if err != nil {
			return "", errors.New("save controller failed")
		}
		buf, err := json.Marshal(controllerMap)
		if err != nil {
			return "", errors.New("save controller failed")
		}
		err = os.WriteFile(controllerPath, buf, 0664)
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
			return "", errors.New("WORKER_ALREADY_IN_CLUSTER:" + worker.Id)
		}
		workerMap[worker.Id] = worker
		fmt.Println(workerMap)
	} else {
		fmt.Printf("err: %s\n", pieces)
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	return meatadata.Node.ClusterId, nil
}
