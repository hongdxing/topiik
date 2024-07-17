/***
**
**
**
**
**/

package executor

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"topiik/internal/consts"
	"topiik/shared"
)

/***
** Get ttl of the the key in seconds, or set ttl of the key if optional seconds provided
** Parameter
**	- pieces: command line that CMD(TTL) stripped, the first piece is the KEY
** Return if optional seconds in command:
**	- The ttl seconds if key exist, and ever set ttl via TTL command or SET command
**	- -1 if the key exist and never set ttl via TTL command or SET command
**	- KEY_NOT_EXIST if key not exist or expired already
** Return if optional seconds in command:
**	- The seconds itself if key exists
**	- KEY_NOT_EXIST if key not exist or expired already
** Syntax: TTL KEY [seconds]
**
**/
func ttl(pieces []string) (ttl int64, err error) {

	if len(pieces) < 1 || len(pieces) > 2 {
		return 0, errors.New(RES_SYNTAX_ERROR)
	}

	key := strings.TrimSpace(pieces[0])
	if val, ok := shared.MemMap[key]; ok {
		if len(pieces) == 1 {
			if val.Exp == consts.INT64_MAX {
				return -1, nil // never expire
			} else if val.Exp-time.Now().UTC().Unix() < 0 {
				// delete the key
				delete(shared.MemMap, key)
				return 0, errors.New(RES_KEY_NOT_EXIST)
			}
			return val.Exp - time.Now().UTC().Unix(), nil
		} else {
			ttl, err = strconv.ParseInt(pieces[1], 10, 64)
			if err != nil {
				return 0, errors.New(RES_SYNTAX_ERROR)
			}
			val.Exp = time.Now().UTC().Unix() + int64(ttl)
			return ttl, nil
		}
	} else {
		return 0, errors.New(RES_KEY_NOT_EXIST)
	}
}
