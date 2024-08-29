/*
* author: Duan HongXing
* date: 28 Aug, 2024
* desc:
*
 */

package keyy

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"topiik/executor/shared"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

/*
* Get ttl of the the key in seconds, or set ttl of the key if optional seconds provided
* Syntax: TTL KEY [seconds] [AT]
*
 */
func Ttl(req datatype.Req) (ttl int64, err error) {
	key := string(req.Keys[0])
	if val, ok := memo.MemMap[key]; ok {
		if strings.TrimSpace(req.Args) == "" {
			if val.Ttl == consts.INT64_MAX {
				/* -1 means never expire */
				return -1, nil
			} else if shared.IsKeyExpired(key, val.Epo, val.Ttl) {
				// delete the key
				//delete(memo.MemMap, key)
				return -2, nil
			}
			return val.Ttl - time.Now().UTC().Unix(), nil
		} else {
			pieces := strings.Split(req.Args, consts.SPACE)
			ttl, err = strconv.ParseInt(pieces[0], 10, 64)
			if err != nil {
				return 0, errors.New(resp.RES_SYNTAX_ERROR)
			}
			val.Ttl = time.Now().UTC().Unix() + int64(ttl)
			return ttl, nil
		}
	} else {
		/* -1 means key not exist */
		return -2, nil
	}
}
