/***
** author: duan hongxing
** date: 26 Jun 2024
** desc:
**
**/

package executor

import (
	"errors"
	"strconv"
	"strings"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/shared"
)

/***
** Increase a KEY, if KEY not exists, create the KEY first
** Parameters
** 	- pieces: command line that CMD stripped, the first piece is the KEY
** Return
**	- The value after increase if success
**	- INVALID_OPT if the KEY is NOT STRING
**
** Syntax: INCR KEY [num]
**/
func INCR(pieces []string) (result string, err error) {
	if (len(pieces)) == 1 { // KEY
		i := 0
		key := strings.TrimSpace(pieces[0])
		i, err = preINCR(key)
		if err != nil {
			return "", err
		}
		i++
		shared.MemMap[key].String = []byte(strconv.Itoa(i))
		return strconv.Itoa(i), nil
	} else if len(pieces) == 2 { // KEY num
		var i int
		var num int
		num, err = strconv.Atoi(pieces[1])
		if err != nil {
			return RES_NIL, errors.New(RES_SYNTAX_ERROR)
		}
		key := strings.TrimSpace(pieces[0])
		i, err = preINCR(key)
		if err != nil {
			return "", err
		}
		i += num
		shared.MemMap[key].String = []byte(strconv.Itoa(i))
		return strconv.Itoa(i), nil
	} else {
		return RES_NIL, errors.New(RES_WRONG_NUMBER_OF_ARGS)
	}
}

func preINCR(key string) (i int, err error) {
	if val, ok := shared.MemMap[key]; ok {
		i, err = strconv.Atoi(string(val.String))
		if err != nil {
			return i, errors.New(RES_DATA_TYPE_NOT_MATCH)
		}
	} else {
		shared.MemMap[key] = &datatype.TValue{
			Type:   datatype.TTYPE_STRING,
			String: []byte("0"),
			Expire: consts.UINT32_MAX}
	}
	return i, nil
}
