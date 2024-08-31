/*
* ©2024 www.topiik.com
* Author: duan hongxing
* Date: 22 Jun 2024
 */

package str

import (
	"errors"
	"topiik/executor/shared"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

/*
* Desc: get STRING value
* Parameter:
*	- pieces: command line that CMD stripped, the first piece is the KEY
* Return:
*	- the value of the key if success
*	- RES_NIL if KEY not exists
*	- RES_DATA_TYPE_MISMATCH if the key found but wrong type
* Syntax:
* 	GET KEY
 */
func Get(req datatype.Req) (result string, err error) {
	key := string(req.Keys[0])
	if val, ok := memo.MemMap[key]; ok {
		if shared.IsKeyExpired(key, val.Ttl) {
			return resp.RES_NIL, nil
		}
		if val.Typ != datatype.V_TYPE_STRING {
			return "", errors.New(resp.RES_DATA_TYPE_MISMATCH)
		}
		return string(val.Str), nil
	} else {
		return resp.RES_NIL, nil
	}
}
