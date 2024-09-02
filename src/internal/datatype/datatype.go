package datatype

/* structs for SHOW command */
type ClusterData struct {
	Controllers []NodeData
	Workers     []NodeData
	Partitions  []PartitionData
}

type NodeData struct {
	Id      string
	Address string
}

type PartitionData struct {
	Id    string
	Slots []SlotData
	Nodes []NodeData
}

type SlotData struct {
	From uint16
	To   uint16
}
