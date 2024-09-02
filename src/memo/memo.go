// ©2024 www.topiik.com
// The key/value map
package memo

import (
	"container/list"
)

const (
	V_TYPE_STRING = 1
	V_TYPE_LIST   = 2
	V_TYPE_HASH   = 3
	V_TYPE_SET    = 4
	V_TYPE_ZSET   = 5
	V_TYPE_GEO    = 6
)

type TValue struct {
	Typ uint8
	Str []byte
	Lst *list.List
	//Hsh
	//Set
	//Zet
	//
	Epo int64 // epoch of created
	Ttl int64 // ttl AT
}

// the kv map
var MemMap = make(map[string]*TValue)
