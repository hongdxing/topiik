// ©2024 www.topiik.com
// author: Duan Hongxing
// data: 3 Jul 2024

package cluster

import (
	"topiik/node"
)

type WorkerGroup struct {
	LeaderNodeId string            // dynamic update
	Nodes        []*node.NodeSlim  // nodes in the group
	Slots        map[uint16]string // partition slots of group
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

// datatype of ControllerInfo and WorkerInfo
type NodesInfo struct {
	ClusterId string
	Nodes     map[string]node.NodeSlim
}

type PartitionInfo struct {
	ClusterId string
	PtnMap    map[string]*node.Partition
	Slots     map[uint16]string
}

// store dynamic node status in runtime
type NodeStatus struct {
	Role                 uint8  // Raft role
	Term                 uint   // Raft term
	Heartbeat            uint16 // Raft heartbeat timeout
	HeartbeatAt          int64  // The UTC milli seconds when heartbeat received from Leader
	LeaderControllerAddr string // Leader Controller updated by append entry, for redirect to Cluster Leader
}
