/*
* author: duan hongxing
* date: 31 Jul, 2024
* desc:
*
 */

package persistence

import "topiik/internal/logger"

/* binlog seq and cmd len */
const preMsgLen int = 12

/* current partition binlog seq */
var binlogSeq int64 = 0

/* log */
var l = logger.Get()
