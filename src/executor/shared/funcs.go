/***
* author: duan hongxing
* date: 15 Jul 2024
* desc:

***/

package shared

import (
	"time"
	"topiik/memo"
)

func IsKeyExpired(key string, exp int64) bool {
	if exp-time.Now().UTC().Unix() < 0 {
		delete(memo.MemMap, key)
		return true
	}
	return false
}
