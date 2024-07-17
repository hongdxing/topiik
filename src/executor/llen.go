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
** Get lenght of the list
** Parameters:
**	- pieces: command line that CMD stripped, the first piece is the KEY
** Return:
**	- Length of the LIST
**	- SYNTAX_ERROR if synctax error
**	- NIL if KEY not exists
** Syntax: LLEN KEY
**/
func llen(pieces []string) (result int, err error) {
	if len(pieces) != 1 {
		return 0, errors.New(RES_SYNTAX_ERROR)
	}
	key := strings.TrimSpace(pieces[0])
	if val, ok := shared.MemMap[key]; ok {
		if val.Typ != datatype.V_TYPE_LIST {
			return 0, errors.New(RES_DATA_TYPE_NOT_MATCH)
		}
		return val.Lst.Len(), nil
	} else {
		return 0, errors.New(RES_NIL)
	}
}
