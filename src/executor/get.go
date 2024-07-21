/***
* author: duan hongxing
* date: 22 Jun 2024
* desc:

***/

package executor

import (
	"errors"
	"strings"
	"topiik/internal/datatype"
	"topiik/memo"
)

/***
** Desc: get STRING value
** Parameter:
**	- pieces: command line that CMD stripped, the first piece is the KEY
** Return:
**	- the value of the key if success
**	- RES_NIL if KEY not exists
**	- RES_DATA_TYPE_NOT_MATCH if the key found but wrong type
** Syntax:
** 	GET KEY
**/
func get(pieces []string) (result string, err error) {
	if len(pieces) != 1 {
		return "", errors.New(RES_SYNTAX_ERROR)
	}

	key := strings.TrimSpace(pieces[0])
	if val, ok := memo.MemMap[key]; ok {
		if isKeyExpired(key, val.Exp) {
			return "", errors.New(RES_KEY_NOT_EXIST)
		}
		if val.Typ != datatype.V_TYPE_STRING {
			return "", errors.New(RES_DATA_TYPE_NOT_MATCH)
		}
		return string(val.Str), nil
	} else {
		return "", errors.New(RES_KEY_NOT_EXIST)
	}
}
