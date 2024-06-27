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
	kv :=make(map[string]string)
	for i := 0; i < len(pieces)-1; i += 2 {
		key := strings.TrimSpace(pieces[i])
		if val, ok := shared.MemMap[key]; ok {
			if val.Type != datatype.TTYPE_STRING {
				return 0, errors.New(RES_DATA_TYPE_NOT_MATCH + ":" + key)
			}
		}
		kv[key] = pieces[i+1]
	}
	for k, v:= range kv{
		shared.MemMap[k]
	}
	return 1, nil
}
