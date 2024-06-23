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
)

/***
** Desc: get STRING value
** Parameter:
**	- pieces: command line that CMD stripped, the first piece is the KEY
** Return:
**	- the value of the key if success
**	- NIL if KEY not exists
**	- invalide operation if the key found but wrong type
** Syntax:
** 	GET KEY
**/
func get(pieces []string) (result string, err error) {
	fmt.Println(pieces)
	if len(pieces) != 1 {
		return "", errors.New(RES_SYNTAX_ERROR)
	}

	if val, ok := memMap[strings.TrimSpace(pieces[0])]; ok {
		if val.Type != datatype.TTYPE_STRING {
			return "", errors.New(RES_INVALID_OP)
		}
		return string(val.String), nil
	} else {
		return "", errors.New(RES_NIL)
	}
}
