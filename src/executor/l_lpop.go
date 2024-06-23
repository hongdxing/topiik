/***
** author: duan hongxing
** date: 23 Jun 2024
** desc:
**
**/

package executor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
** Command: LPOP KEY [COUNT]
**/
func lPop(args []string) (result []string, err error) {
	count := 1
	if len(args) > 2 {
		return nil, errors.New(RES_WRONG_NUMBER_OF_ARGS)
	} else if len(args) == 2 {
		count, err = strconv.Atoi(args[2])
		if err != nil {
			return nil, errors.New(RES_WRONG_ARG)
		}
		if count < 1 {
			return nil, errors.New(RES_WRONG_ARG)
		}
	}
	key := strings.TrimSpace(args[0])
	fmt.Printf("%s", key)
	fmt.Println(memMap)
	if val, ok := memMap[key]; ok {
		fmt.Println(val)
		if val.Type != datatype.TTYPE_LIST {
			return result, errors.New(RES_DATA_TYPE_NOT_MATCH)
		}
		for i := 0; i < count; i++ {
			if ele := val.List.Front(); ele != nil {
				//ele = ele.Next()
				result = append(result, ele.Value.(string))
				val.List.Remove(ele)
			} else {
				break
			}
		}
		return result, nil
	}
	return result, errors.New(RES_NIL)
}
