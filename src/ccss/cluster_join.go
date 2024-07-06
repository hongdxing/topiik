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
**	id: the salor id
**	host:port: the salor listen host:port
**
**/
func clusterJoin(pieces []string) (result string, err error) {
	if len(pieces) != 3 {
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	salor := Salor{
		Id:      pieces[1],
		Address: pieces[2],
	}
	if _, ok := salorMap[salor.Id]; ok {
		return "", errors.New("SALOR_ALREADY_IN_CLUSTER:" + salor.Id)
	}
	salorMap[salor.Id] = salor
	fmt.Println(salorMap)
	return "OK", nil
}
