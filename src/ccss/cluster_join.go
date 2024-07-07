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
)

/***
**
** Parameter:
**	- pieces: JON_ACK id host:port
**
** Syntax: CLUSTER JOIN_ACK id host:port
**	id: the worker id
**	host:port: the worker listen host:port
**
**/
func clusterJoin(pieces []string) (result string, err error) {
	if len(pieces) != 3 {
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	worker := Worker{
		Id:      pieces[1],
		Address: pieces[2],
	}
	if _, ok := workerMap[worker.Id]; ok {
		return "", errors.New("WORKER_ALREADY_IN_CLUSTER:" + worker.Id)
	}
	workerMap[worker.Id] = worker
	fmt.Println(workerMap)
	return "OK", nil
}
