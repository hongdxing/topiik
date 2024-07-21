/***
** author: duan hongxing
** date: 23 Jun 2024
** desc:
**
**/

package executor

import (
	"container/list"
	"errors"
	"topiik/internal/command"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/memo"
)

/***
** Push element(s) to left(header) of list
** Parameter:
** 	- args: the arguments, command line that CMD(LPUSH) stripped, the first piece is the KEY
** Return:
** 	- Lenght of the list after push
**	- INVALID_OP if the key exists but data type is not list
**
** Syntax: LPUSH|RPUSH key value1 [... valueN]
**/
func pushList(pieces []string, icmd int16) (result int, err error) {
	if len(pieces) < 2 { // except KEY, at least need one value
		return 0, errors.New(RES_WRONG_NUMBER_OF_ARGS)
	}
	/*pieces, err := util.SplitCommandLine(args[1])
	if err != nil {
		return 0, err
	}*/
	key := pieces[0]
	if memo.MemMap[key] == nil {
		memo.MemMap[key] = &datatype.TValue{
			Typ: datatype.V_TYPE_LIST,
			Lst: list.New(),
			Exp: consts.UINT32_MAX,
		}
	}
	if icmd == command.LPUSH {
		for _, piece := range pieces[1:] {
			memo.MemMap[key].Lst.PushFront(piece)
		}
	} else if icmd == command.LPUSHR {
		for _, piece := range pieces[1:] {
			memo.MemMap[key].Lst.PushBack(piece)
		}
	} else {
		return 0, errors.New(consts.RES_INVALID_CMD)
	}

	return memo.MemMap[key].Lst.Len(), nil
}
