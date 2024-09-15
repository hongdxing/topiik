// author: Duan Hongxing
// date: 24 Aug, 2024

package node

type Node struct {
	Id        string
	ClusterId string
	GroupId   string
	Role      string // controller or worker
	Addr      string
	Addr2     string
}

type NodeSlim struct {
	Id    string
	Addr  string
	Addr2 string
}

type Slot struct {
	From uint16
	To   uint16
}

type Partition struct {
	Id           string               // partition id
	LeaderNodeId string               // partition leader node id
	NodeSet      map[string]*NodeSlim // member nodes
	//SlotFrom     uint16
	//SlotTo       uint16
	Status string
}
