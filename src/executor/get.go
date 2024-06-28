/***
* author: duan hongxing
* date: 22 Jun 2024
* desc:

***/

package executor

import (
	"errors"
	"fmt"
	"strings"
	"topiik/internal/datatype"
	"topiik/shared"
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
	fmt.Println(pieces)
	if len(pieces) != 1 {
		return "", errors.New(RES_SYNTAX_ERROR)
	}

	if val, ok := shared.MemMap[strings.TrimSpace(pieces[0])]; ok {
		if val.Typ != datatype.V_TYPE_STRING {
			return "", errors.New(RES_DATA_TYPE_NOT_MATCH)
		}
		return string(val.Str), nil
	} else {
		return "", errors.New(RES_KEY_NOT_EXIST)
	}
}
