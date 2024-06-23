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
	"topiik/internal/consts"
	"topiik/internal/datatype"
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
		fmt.Println(memMap)
		memMap[strings.TrimSpace(pieces[0])] = &datatype.TValue{
			Type:   datatype.TTYPE_STRING,
			String: []byte(pieces[1]),
			Expire: consts.UINT32_MAX}
		return RES_OK, nil
	} else {
		return "", errors.New(RES_SYNTAX_ERROR)
	}
}
