/***
** author: duan hongxing
** data: 4 Jul 2024
** desc:
**
**/

package ccss

import (
	"errors"
	"fmt"
	"strings"
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
		controller := Controller{
			Id:      pieces[1],
			Address: pieces[2],
		}
		if _, ok := controllerMap[controller.Id]; ok {
			return "", errors.New("WORKER_ALREADY_IN_CLUSTER:" + controller.Id)
		}
		controllerMap[controller.Id] = controller
		fmt.Println(controllerMap)
	} else if strings.ToUpper(pieces[3]) == ROLE_WORKER {
		worker := Worker{
			Id:      pieces[1],
			Address: pieces[2],
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
