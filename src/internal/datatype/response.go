/***
**
**
**
**/

package datatype

type Response[T any] struct {
	R bool
	M T
}

// cluster
type PartitionInfo struct {
	Id string
}

type NodeInfo struct {
	Id         string
	Partitions []PartitionInfo
	Replicas   []PartitionInfo
}

type ClusterInfoResponse struct {
	Nodes []NodeInfo
}
