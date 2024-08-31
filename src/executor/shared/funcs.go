/*
* Author: Duan Hongxing
* Date: 15 Jul, 2024
* Desc:
*
 */

package shared

import (
	"topiik/internal/consts"
	"topiik/internal/util"
	"topiik/memo"
)

func IsKeyExpired(key string, ttl int64) bool {
	if ttl == consts.INT64_MIN {
		return false
	}
	if ttl < util.GetUtcEpoch() {
		delete(memo.MemMap, key)
		return true
	}
	return false
}
