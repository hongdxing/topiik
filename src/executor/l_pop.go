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
	"topiik/internal/datatype"
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
func popList(args []string, cmd string) (result []string, err error) {
	count := 1
	if len(args) > 2 {
		return nil, errors.New(RES_WRONG_NUMBER_OF_ARGS)
	} else if len(args) == 2 {
		count, err = strconv.Atoi(args[1])
		if err != nil {
			return nil, errors.New(RES_WRONG_ARG)
		}
		if count < 1 {
			return nil, errors.New(RES_WRONG_ARG)
		}
	}
	key := strings.TrimSpace(args[0])
	if val, ok := memMap[key]; ok {
		if val.Type != datatype.TTYPE_LIST {
			return result, errors.New(RES_DATA_TYPE_NOT_MATCH)
		}

		var eleToBeRemoved []*list.Element
		if cmd == command.LPOP { //LPOP
			looper := 0
			for ele := val.TList.Front(); ele != nil && looper < count; ele = ele.Next() {
				looper++
				result = append(result, ele.Value.(string))
				eleToBeRemoved = append(eleToBeRemoved, ele)
			}
		} else if cmd == command.RPOP { //RPOP
			looper := 0
			for ele := val.TList.Back(); ele != nil && looper < count; ele = ele.Prev() {
				looper++
				result = append(result, ele.Value.(string))
				eleToBeRemoved = append(eleToBeRemoved, ele)
			}
		} else {
			return result, errors.New(RES_INVALID_CMD) // This should never happen
		}
		// remove ele from list
		for _, ele := range eleToBeRemoved {
			val.TList.Remove(ele)
		}
		// if no element, delete the list
		if val.TList.Len() == 0 {
			delete(memMap, key)
		}

		return result, nil
	}
	return result, errors.New(RES_NIL)
}
