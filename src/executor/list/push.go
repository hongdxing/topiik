/*
* author: duan hongxing
* date: 23 Jun 2024
* desc:
*
 */

package list

import (
	"container/list"
	"errors"
	"topiik/internal/command"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/memo"
)

/*
* Push element(s) to left(header) of list
* Parameter:
* 	- args: the arguments, command line that CMD(LPUSH) stripped, the first piece is the KEY
* Return:
* 	- Lenght of the list after push
*	- INVALID_OP if the key exists but data type is not list
*
* Syntax: LPUSH|RPUSH key value1 [... valueN]
 */
func Push(req datatype.Req, icmd uint8) (result int, err error) {
	key := string(req.Keys[0])
	if memo.MemMap[key] == nil {
		memo.MemMap[key] = &datatype.TValue{
			Typ: datatype.V_TYPE_LIST,
			Lst: list.New(),
			Exp: consts.UINT32_MAX,
		}
	}
	if icmd == command.LPUSH_I {
		for _, val := range req.Vals {
			memo.MemMap[key].Lst.PushFront(val)
		}
	} else if icmd == command.LPUSHR_I {
		for _, val := range req.Vals {
			memo.MemMap[key].Lst.PushBack(val)
		}
	} else {
		return 0, errors.New(consts.RES_INVALID_CMD)
	}

	return memo.MemMap[key].Lst.Len(), nil
}
