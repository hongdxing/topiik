/***
* author: duan hongxing
* date: 22 Jun 2024
* desc:
***/

package executor

import (
	"errors"
	"strings"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/shared"
)

/***
** Set STRING kv
** Parameter:
**	- pieces: command line that CMD(SET) stripped, the first piece is the KEY
** Return:
**	- OK if success
**	- SYNTAX_ERROR if syntax error
**/
func set(pieces []string) (result string, err error) {
	if len(pieces) == 2 {
		shared.MemMap[strings.TrimSpace(pieces[0])] = &datatype.TValue{
			Typ: datatype.V_TYPE_STRING,
			Str: []byte(pieces[1]),
			Exp: consts.UINT32_MAX}
		return RES_OK, nil
	} else {
		return "", errors.New(RES_SYNTAX_ERROR)
	}
}
