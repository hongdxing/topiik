//author: Duan HongXing
//date: 28 Aug, 2024

package keyy

import (
	"topiik/internal/datatype"
	"topiik/memo"
)

// Chech if key(s) exists
// Return array of T if exist or F if not exist
func Exists(req datatype.Req) (rslt []string, err error) {
	for _, keyB := range req.Keys {
		key := string(keyB)
		if _, ok := memo.MemMap[string(key)]; ok {
			rslt = append(rslt, "T")
		} else {
			rslt = append(rslt, "F")
		}
	}
	return rslt, err
}
