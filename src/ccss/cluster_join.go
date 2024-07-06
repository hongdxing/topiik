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
**	id: the sailor id
**	host:port: the sailor listen host:port
**
**/
func clusterJoin(pieces []string) (result string, err error) {
	if len(pieces) != 3 {
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	sailor := Sailor{
		Id:      pieces[1],
		Address: pieces[2],
	}
	if _, ok := sailorMap[sailor.Id]; ok {
		return "", errors.New("SAILOR_ALREADY_IN_CLUSTER:" + sailor.Id)
	}
	sailorMap[sailor.Id] = sailor
	fmt.Println(sailorMap)
	return "OK", nil
}
