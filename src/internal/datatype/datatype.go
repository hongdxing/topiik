package datatype

import "container/list"

type Key struct {
	TheKey string
	Expire int
}

const (
	TTYPE_STRING = 1
	TTYPE_LIST   = 2
	TTYPE_HASH   = 3
	TTYPE_SET    = 4
	TTYPE_ZSET   = 5
	TTYPE_GEO    = 6
)

type TValue struct {
	Type   uint8
	String []byte
	TList  *list.List
	/***
	* unint32 Max(-1): no expire
	* else: seconds to epxire
	* max value: 4294967295 = Sunday, February 7, 2106 6:28:15 AM
	 */
	Expire uint32
}
