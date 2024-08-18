/***
* author: duan hongxing
* date: 22 Jun 2024
* desc:

***/

package str

import (
	"errors"
	"strings"
	"topiik/executor/shared"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
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
func Get(req datatype.Req) (result string, err error) {

	/*if len(pieces) != 1 {
		return "", errors.New(resp.RES_SYNTAX_ERROR)
	}*/

	key := strings.TrimSpace(req.KEYS[0])
	if val, ok := memo.MemMap[key]; ok {
		if shared.IsKeyExpired(key, val.Exp) {
			return "", errors.New(resp.RES_KEY_NOT_EXIST)
		}
		if val.Typ != datatype.V_TYPE_STRING {
			return "", errors.New(resp.RES_DATA_TYPE_NOT_MATCH)
		}
		return string(val.Str), nil
	} else {
		return "", errors.New(resp.RES_KEY_NOT_EXIST)
	}
}
