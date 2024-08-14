package node

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

type Slot struct {
	From uint16
	To   uint16
}

type Partition struct {
	Id           string              // Id of the partition, random 16 alphnum
	LeaderNodeId string              // The Node Id where the Leader Partition located
	NodeSet      map[string]NodeSlim // the byte value is not important
	Slots        []Slot
}
