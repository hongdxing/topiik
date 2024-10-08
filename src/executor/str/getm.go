/***
** author: duan hongxing
** date: 26 Jun 2024
** desc:
**
**/

package str

import (
	"errors"
	"topiik/executor/shared"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

/*
* Get multi values of STRING KEYs
* Parameters:
*	- req
* Return:
*	- list of value, the length of the returned values is the same as lenght of KEYs, if some key not exist then NIL in that position
*	- RES_DATA_TYPE_MISMATCH if any KEY is not STRING
*
* Syntax: GETM KEY1 KEY2 [... KEYn]
 */
func GetM(req datatype.Req) (result []string, err error) {
	if len(req.Keys) < 1 {
		return nil, errors.New(resp.RES_SYNTAX_ERROR)
	}
	for _, key := range req.Keys {
		skey := string(key)
		if val, ok := memo.MemMap[skey]; ok {
			if shared.IsKeyExpired(skey, val.Ttl) {
				result = append(result, "")
			}
			if val.Typ != memo.V_TYPE_STRING {
				return nil, errors.New(resp.RES_DATA_TYPE_MISMATCH)
			}
			result = append(result, string(val.Str))
		} else {
			result = append(result, "")
		}
	}
	return result, nil
}
