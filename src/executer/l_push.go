/***
** author: duan hongxing
** date: 23 Jun 2024
** desc:
**
**/

package executer

import (
	"container/list"
	"errors"
	"topiik/internal/command"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/util"
	"topiik/shared"
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
func pushList(args []string, CMD string) (result int, err error) {
	if len(args) < 2 { // except KEY, at least need one value
		return 0, errors.New(RES_WRONG_NUMBER_OF_ARGS)
	}
	pieces, err := util.SplitCommandLine(args[1])
	if err != nil {
		return 0, err
	}
	key := args[0]
	if shared.MemMap[key] == nil {
		shared.MemMap[key] = &datatype.TValue{
			Typ:   datatype.V_TYPE_LIST,
			Lst:  list.New(),
			Exp: consts.UINT32_MAX,
		}
	}
	if CMD == command.LPUSH {
		for _, piece := range pieces {
			shared.MemMap[key].Lst.PushFront(piece)
		}
	} else if CMD == command.LPUSHR {
		for _, piece := range pieces {
			shared.MemMap[key].Lst.PushBack(piece)
		}
	} else {
		return 0, errors.New(RES_INVALID_CMD)
	}

	return shared.MemMap[key].Lst.Len(), nil
}
