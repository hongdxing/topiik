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
	"strconv"
	"strings"
	"topiik/internal/command"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/memo"
)

/***
** Pop element(s) from left(header) of list
** Parameters:
**	- args: the arguments, command line that CMD(LPOP) stripped, the first piece is the KEY
** Return:
**	- list of pop results if success
** 	- err:
**		-
** Syntax: LPOP|RPOP key [COUNT]
**/
func popList(pieces []string, icmd uint8) (result []string, err error) {
	count := 1
	if len(pieces) == 1 {
		//
	} else if len(pieces) == 2 {
		count, err = strconv.Atoi(pieces[1])
		if err != nil {
			return nil, errors.New(RES_WRONG_ARG)
		}
		if count < 1 {
			return nil, errors.New(RES_WRONG_ARG)
		}
	} else {
		return nil, errors.New(RES_SYNTAX_ERROR)
	}
	key := strings.TrimSpace(pieces[0])
	if val, ok := memo.MemMap[key]; ok {
		if val.Typ != datatype.V_TYPE_LIST {
			return result, errors.New(RES_DATA_TYPE_NOT_MATCH)
		}

		var eleToBeRemoved []*list.Element
		if icmd == command.LPOP_I { //LPOP
			looper := 0
			for ele := val.Lst.Front(); ele != nil && looper < count; ele = ele.Next() {
				looper++
				result = append(result, ele.Value.(string))
				eleToBeRemoved = append(eleToBeRemoved, ele)
			}
		} else if icmd == command.LPOPR_I { //RPOP
			looper := 0
			for ele := val.Lst.Back(); ele != nil && looper < count; ele = ele.Prev() {
				looper++
				result = append(result, ele.Value.(string))
				eleToBeRemoved = append(eleToBeRemoved, ele)
			}
		} else {
			return result, errors.New(consts.RES_INVALID_CMD) // This should never happen
		}
		// remove ele from list
		for _, ele := range eleToBeRemoved {
			val.Lst.Remove(ele)
		}
		// if no element, delete the list
		if val.Lst.Len() == 0 {
			delete(memo.MemMap, key)
		}

		return result, nil
	}
	return result, errors.New(RES_KEY_NOT_EXIST)
}
