/*
* author: Duan HongXing
* date: 28 Aug, 2024
* desc:
 */

package keyy

import (
	"topiik/internal/datatype"
	"topiik/memo"
)

func Exists(req datatype.Req) (rslt int64, err error) {
	for _, keyB := range req.Keys {
		key := string(keyB)
		if _, ok := memo.MemMap[string(key)]; ok {
			rslt++
		}
	}
	return rslt, err
}