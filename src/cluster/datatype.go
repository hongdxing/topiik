/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package cluster

type Metadata struct {
	Node Node
}

type Cluster struct {
	Id          string // auto generated when INIT a cluster
	Partitions  uint16 // number of partitions, given parameter when INIT a cluster
	Replicas    uint16 // number of replicas, recommend 3 relicas at most, given parameter when INIT a cluster
	Ver         uint   // compare which is more lastest
	Controllers map[string]NodeSlim
	Workers     map[string]NodeSlim
}

type Node struct {
	Id        string
	ClusterId string
	Address   string
	Address2  string
}

type NodeSlim struct {
	Id       string
	Address  string
	Address2 string
}

type Partition struct {
	Id           string   // Id of the partition, random 16 alphnum
	LeaderNodeId string   // The Node Id where the Leader Partition located
	NodeIds      []string // The Nodes where the Partition located
}
