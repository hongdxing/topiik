// ©2024 www.topiik.com
// author: Duan Hongxing
// data: 3 Jul, 2024

package cluster

import (
	"topiik/node"
)

type PartitionInfo struct {
	ClusterId string
	Ptns      map[string]*Partition
}

type Partition struct {
	Id           string
	LeaderNodeId string                   // dynamic update
	Nodes        map[string]node.NodeSlim // nodes in the group
	Slots        map[uint16]bool          // partition slots of group, the bool value is not important
}

type Persistors struct {
	LeaderNodeId string
	Nodes        []*node.NodeSlim
}

type Cluster struct {
	Id   string // auto generated when INIT a cluster
	Ver  uint   // compare which is more lastest
	Ctls map[string]node.NodeSlim
	Wkrs map[string]node.NodeSlim
}

// store dynamic node status in runtime
type NodeStatus struct {
	RaftRole             uint8  // Raft role
	Term                 uint   // Raft term
	Heartbeat            uint16 // Raft heartbeat timeout
	HeartbeatAt          int64  // The UTC milli seconds when heartbeat received from Leader
	LeaderControllerAddr string // Leader Controller updated by append entry, for redirect to Cluster Leader
}
