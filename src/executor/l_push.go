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
	"strings"
	"topiik/internal/command"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/util"
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
	key := strings.Clone(args[0])
	if memMap[key] == nil {
		memMap[key] = &datatype.TValue{
			Type:   datatype.TTYPE_LIST,
			TList:  list.New(),
			Expire: consts.UINT32_MAX,
		}
	}
	if CMD == command.LPUSH {
		for _, piece := range pieces {
			memMap[key].TList.PushFront(piece)
		}
	} else if CMD == command.LPUSHR {
		for _, piece := range pieces {
			memMap[key].TList.PushBack(piece)
		}
	} else {
		return 0, errors.New(RES_INVALID_CMD)
	}

	return memMap[key].TList.Len(), nil
}
