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
	"topiik/memo"
)

/***
** Get multi values of STRING KEYs
** Parameters:
**	- pieces: command line that CMD stripped, the first piece is the KEY
** Return:
**	- list of value, the length of the returned values is the same as lenght of KEYs, if some key not exist then NIL in that position
**	- RES_DATA_TYPE_NOT_MATCH if any KEY is not STRING
**
** Syntax: GETM KEY1 KEY2 [... KEYn]
**/
func getM(pieces []string) (result []string, err error) {
	if len(pieces) < 1 {
		return nil, errors.New(RES_SYNTAX_ERROR)
	}
	for _, key := range pieces {
		key = strings.TrimSpace(key)
		if val, ok := memo.MemMap[key]; ok {
			if isKeyExpired(key, val.Exp) {
				result = append(result, "")
			}
			if val.Typ != datatype.V_TYPE_STRING {
				return nil, errors.New(RES_DATA_TYPE_NOT_MATCH)
			}
			result = append(result, string(val.Str))
		} else {
			result = append(result, "")
		}
	}
	return result, nil
}
