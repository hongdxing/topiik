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
	"strconv"
	"topiik/internal/command"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

/*
* Pop element(s) from left(header) of list
* Parameters:
*	- args: the arguments, command line that CMD(LPOP) stripped, the first piece is the KEY
* Return:
*	- list of pop results if success
* 	- err:
*		-
* Syntax: LPOP|RPOP key [COUNT]
 */
func Pop(req datatype.Req, icmd uint8) (rslt []string, err error) {
	count := 1
	if len(req.Args) != 0 { // if args has value, the only value should be the count
		count, err = strconv.Atoi(req.Args)
		if err != nil {
			return nil, errors.New(resp.RES_WRONG_ARG)
		}
		if count < 1 {
			return nil, errors.New(resp.RES_WRONG_ARG)
		}
	}
	key := string(req.Keys[0])
	if val, ok := memo.MemMap[key]; ok {
		if val.Typ != memo.V_TYPE_LIST {
			return rslt, errors.New(resp.RES_DATA_TYPE_MISMATCH)
		}

		var eleToBeRemoved []*list.Element
		if icmd == command.LPOP_I { //LPOP
			looper := 0
			for ele := val.Lst.Front(); ele != nil && looper < count; ele = ele.Next() {
				looper++
				rslt = append(rslt, string(ele.Value.([]byte)))
				eleToBeRemoved = append(eleToBeRemoved, ele)
			}
		} else if icmd == command.LPOPR_I { //RPOP
			looper := 0
			for ele := val.Lst.Back(); ele != nil && looper < count; ele = ele.Prev() {
				looper++
				rslt = append(rslt, string(ele.Value.([]byte)))
				eleToBeRemoved = append(eleToBeRemoved, ele)
			}
		} else {
			return rslt, errors.New(consts.RES_INVALID_CMD) // This should never happen
		}
		// remove ele from list
		for _, ele := range eleToBeRemoved {
			val.Lst.Remove(ele)
		}
		// if no element, delete the list
		if val.Lst.Len() == 0 {
			delete(memo.MemMap, key)
		}

		return rslt, nil
	}
	return rslt, errors.New(resp.RES_KEY_NOT_EXIST)
}
