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
	"topiik/executor/shared"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/util"
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
		// get ttl
		if strings.TrimSpace(req.Args) == "" {
			if val.Ttl == consts.INT64_MIN {
				/* -1 means never expire */
				return -1, nil
			} else if shared.IsKeyExpired(key, val.Ttl) {
				// delete the key
				//delete(memo.MemMap, key)
				return -2, nil
			}
			ttl = val.Ttl - util.GetUtcEpoch()
			// only when ttl is INT64_MIN, should return -1
			// if ttl is lt 0, then return -2
			if ttl <= 0 {
				ttl = -2
			}
			return ttl, nil
		} else {
			// set ttl
			pieces := strings.Split(req.Args, consts.SPACE)
			ttl, err = strconv.ParseInt(pieces[0], 10, 64)
			if err != nil {
				return 0, errors.New(resp.RES_SYNTAX_ERROR)
			}
			// convert to epoch
			val.Ttl = ttl + util.GetUtcEpoch()
			return ttl, nil
		}
	} else {
		/* -1 means key not exist */
		return -2, nil
	}
}
