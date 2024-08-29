/*
* Author: Duan Hongxing
* Date: 15 Jul, 2024
* Desc:
*
 */

package shared

import (
	"time"
	"topiik/internal/consts"
	"topiik/memo"
)

func IsKeyExpired(key string, epo int64, ttl int64) bool {
	if ttl == consts.INT64_MAX {
		return false
	}
	if (epo + ttl) < time.Now().UTC().Unix() {
		delete(memo.MemMap, key)
		return true
	}
	return false
}
