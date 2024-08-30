package datatype

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
	Ttl int64 // ttl
}

/* structs for SHOW command */
type Cluster struct {
	Controllers []Node
	Workers     []Node
	Partitions  map[string][]Partition
}

type Node struct {
	Id      string
	Address string
}

type Partition struct {
	Id    string
	Slots []int
}

type Slot struct {
	From uint16
	To   uint16
}
