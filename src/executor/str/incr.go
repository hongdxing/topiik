/***
** author: duan hongxing
** date: 26 Jun 2024
** desc:
**
**/

package str

import (
	"errors"
	"strconv"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

/***
** Increase a KEY, if KEY not exists, create the KEY first
** Parameters
** 	- req
** Return
**	- The value after increase if success
**	- INVALID_OPT if the KEY is NOT STRING
**
** Syntax: INCR KEY [num]
**/
func Incr(req datatype.Req) (result int64, err error) {
	if len(req.KEYS) == 0 {
		return 0, errors.New(resp.RES_SYNTAX_ERROR)
	}
	if req.ARGS == "" { // KEY
		var i int64 = 0
		key := string(req.KEYS[0])
		i, err = preINCR(key)
		if err != nil {
			return 0, err
		}
		i++
		memo.MemMap[key].Str = []byte(string(i))
		return i, nil
	} else { // KEY num
		var i int64
		var num int
		num, err = strconv.Atoi(req.ARGS)
		if err != nil {
			return 0, errors.New(resp.RES_SYNTAX_ERROR)
		}
		key := string(req.KEYS[0])
		i, err = preINCR(key)
		if err != nil {
			return 0, err
		}
		i += int64(num)
		memo.MemMap[key].Str = []byte(string(i))
		return i, nil
	}
}

func preINCR(key string) (i int64, err error) {
	if val, ok := memo.MemMap[key]; ok {
		i, err = strconv.ParseInt(string(val.Str), 10, 0)
		if err != nil {
			return i, errors.New(resp.RES_DATA_TYPE_NOT_MATCH)
		}
	} else {
		memo.MemMap[key] = &datatype.TValue{
			Typ: datatype.V_TYPE_STRING,
			Str: []byte("0"),
			Exp: consts.UINT32_MAX}
	}
	return i, nil
}
