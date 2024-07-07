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

type Node struct {
	Id        string
	ClusterId string
}

type Controller struct {
	Id       string
	Address  string
	Address2 string
}

type Worker struct {
	Id       string
	Address  string
	Address2 string
}

type Partition struct {
	Id             string   // Id of the partition, random 16 alphnum
	LeaderWorkerId string   // The Worker Id where the Leader Partition  located
	WorkerIds      []string // The Workers where the Partition located
}
