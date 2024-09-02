package datatype

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
