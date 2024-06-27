/***
** author: duan hongxing
** date: 26 Jun 2024
** desc:
**
**/

package executor

import (
	"errors"
	"strings"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/shared"
)

/***
** Set multi KEY/VALUE
** Parameters:
**	- pieces: command line that CMD stripped
** Return:
**	- number of key set if success
** Syntax: SETM KEY1 VALUE1 [... KEYn VALUEn]
**/
func setM(pieces []string) (result int, err error) {
	if len(pieces)%2 == 1 {
		return 0, errors.New(RES_SYNTAX_ERROR)
	}
	kv := make(map[string]string)
	for i := 0; i < len(pieces)-1; i += 2 {
		key := strings.TrimSpace(pieces[i])
		if val, ok := shared.MemMap[key]; ok {
			if val.Typ != datatype.V_TYPE_STRING {
				return 0, errors.New(RES_DATA_TYPE_NOT_MATCH + ":" + key)
			}
		}
		kv[key] = pieces[i+1]
	}
	for k, v := range kv {
		if val, ok := shared.MemMap[k]; ok {
			val.Str = []byte(v)
		} else {
			shared.MemMap[k] = &datatype.TValue{
				Typ: datatype.V_TYPE_STRING,
				Str: []byte(v),
				Exp: consts.UINT32_MAX,
			}
		}
	}
	return len(kv), nil
}
