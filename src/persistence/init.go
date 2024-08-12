/*
* author: duan hongxing
* date: 31 Jul, 2024
* desc:
*
 */

package persistence

import "topiik/logger"

// current partition binlog seq
var binLogSeq int64 = 0

// partition follower address2 list
var ptnFlrAddr2Lst []string

// log
var l = logger.Get()
