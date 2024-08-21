/***
** author: duan hongxing
** date: 26 Jun 2024
** desc:
**
**/

package str

import (
	"errors"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

/***
** Set multi KEY/VALUE
** Parameters:
**	- pieces: command line that CMD stripped
** Return:
**	- number of key set if success
** Syntax: SETM KEY1 VALUE1 [... KEYn VALUEn]
**/
func SetM(req datatype.Req) (result int, err error) {
	if len(req.KEYS) != len(req.VALS) || len(req.KEYS) == 0 {
		return 0, errors.New(resp.RES_SYNTAX_ERROR)
	}
	kv := make(map[string][]byte)
	for i := 0; i < len(req.KEYS); i++ {
		key := string(req.KEYS[i])
		if val, ok := memo.MemMap[key]; ok { // if the key exists, but not String type, then error
			if val.Typ != datatype.V_TYPE_STRING {
				return 0, errors.New(resp.RES_DATA_TYPE_NOT_MATCH + ":" + key)
			}
		}
		kv[key] = req.VALS[i]
	}
	for k, v := range kv {
		if val, ok := memo.MemMap[k]; ok {
			val.Str = []byte(v)
		} else {
			memo.MemMap[k] = &datatype.TValue{
				Typ: datatype.V_TYPE_STRING,
				Str: []byte(v),
				Exp: consts.UINT32_MAX,
			}
		}
	}
	return len(kv), nil
}
