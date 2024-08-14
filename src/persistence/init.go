/*
* author: duan hongxing
* date: 31 Jul, 2024
* desc:
*
 */

package persistence

import (
	"topiik/logger"
)

// current partition binlog seq
var binLogSeq int64 = 0

// log
var l = logger.Get()
