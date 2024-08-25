/*
* author: duan hongxing
* date: 26 Jun 2024
* desc:
*
 */

package list

import (
	"errors"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

/*
* Get lenght of the list
* Parameters:
*	- pieces: command line that CMD stripped, the first piece is the KEY
* Return:
*	- Length of the LIST
*	- SYNTAX_ERROR if synctax error
*	- ERROR if KEY not exists
* Syntax: LLEN KEY
 */
func Len(req datatype.Req) (result int, err error) {
	key := string(req.Keys[0])
	if val, ok := memo.MemMap[key]; ok {
		if val.Typ != datatype.V_TYPE_LIST {
			return 0, errors.New(resp.RES_DATA_TYPE_NOT_MATCH)
		}
		return val.Lst.Len(), nil
	} else {
		return 0, errors.New(resp.RES_KEY_NOT_EXIST)
	}
}
