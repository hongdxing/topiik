/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package cluster

type Cluster struct {
	Id   string // auto generated when INIT a cluster
	Ptns uint16 // number of partitions, given parameter when INIT a cluster
	Rpls uint16 // number of replicas, recommend 3 relicas at most, given parameter when INIT a cluster
	Ver  uint   // compare which is more lastest
	Ctls map[string]NodeSlim
	Wkrs map[string]Worker
}

type Node struct {
	Id        string
	ClusterId string
	Addr      string
	Addr2     string
}

type NodeSlim struct {
	Id    string
	Addr  string
	Addr2 string
}

type Worker struct {
	Id       string
	Addr     string
	Addr2    string
	LeaderId string // Partition Leader Node Id
	Slots    []Slot // Slots of current Node
}

type Slot struct {
	From uint16
	To   uint16
}

type Partition struct {
	Id           string   // Id of the partition, random 16 alphnum
	LeaderNodeId string   // The Node Id where the Leader Partition located
	NodeIds      []string // The Nodes where the Partition located
}

// store dynamic node status in runtime
type NodeStatus struct {
	Role                 uint8  // Raft role
	Term                 uint   // Raft term
	Heartbeat            uint16 // Raft heartbeat timeout
	HeartbeatAt          int64  // The UTC milli seconds when heartbeat received from Leader
	LeaderControllerAddr string // Leader Controller updated by append entry
}
