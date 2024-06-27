package datatype

import (
	"container/list"
)

type Key struct {
	TheKey string
	Expire int
}

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
	/***
	* unint32 Max(-1): no expire
	* else: seconds to epxire
	* max value: 4294967295 = Sunday, February 7, 2106 6:28:15 AM
	 */
	Exp uint32
}

type ValueT[T any] struct {
	Type  uint8
	Value T
	/***
	* unint32 Max(-1): no expire
	* else: seconds to epxire
	* max value: 4294967295 = Sunday, February 7, 2106 6:28:15 AM
	 */
	Expire uint32
}
