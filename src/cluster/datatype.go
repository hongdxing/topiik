/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package cluster

import (
	"topiik/node"
)

type Cluster struct {
	Id   string // auto generated when INIT a cluster
	Ver  uint   // compare which is more lastest
	Ctls map[string]node.NodeSlim
	Wkrs map[string]Worker
}

type Worker struct {
	Id    string
	Addr  string
	Addr2 string
}

type PartitionInfo struct {
	Ptns   uint16 // number of partitions, given parameter when INIT a cluster
	Rpls   uint16 // number of replicas, recommend 3 relicas at most, given parameter when INIT a cluster
	PtnMap map[string]node.Partition
}

// store dynamic node status in runtime
type NodeStatus struct {
	Role                 uint8  // Raft role
	Term                 uint   // Raft term
	Heartbeat            uint16 // Raft heartbeat timeout
	HeartbeatAt          int64  // The UTC milli seconds when heartbeat received from Leader
	LeaderControllerAddr string // Leader Controller updated by append entry, for redirect to Cluster Leader
}
