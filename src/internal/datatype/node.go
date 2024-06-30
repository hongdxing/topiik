/***
** author: duan hongxing
** date: 30 Jun 2024
** desc:
** There are 3 cluster modes:
** 	1) Single node mode, also non cluster mode, only for development or test purpose
		- One node
		- No partition and replication
**	2) 3 master nodes cluster, but no slaves
		- 3 master nodes
		- The data will be split to 3 partitions, and each partition have 2 replication
		  so, there are 6 replications in total spread evenly in 3 nodes
		- Each node will have 2 replications of different partitions
**	3) 3 master nodes cluster, with 3 more slaves
		- 3 master nodes, 3 slaves nodes
		- The data will be split to 3 partitions, and each partition have 2 replication
		  so, there are 6 replications in total
		- The different with mode 2) is, each master node have a dedicated slave node to store replication
		  so, the 6 partitions will be stored in 6 nodes
**/

package datatype

type Partition struct {
	id       uint8 // 1, 2, 3
	isLeader bool
}

/**
** Under cluster model, each node have 2 partitions
**
**
**/
type Node struct {
	Partitions [2]Partition
}

/***
** Cluster model:
** 1) Single node mode
** 2) 3 masters without slave
** 3) 3 masters with 3 slaves
**/
type Cluster struct {
	Masters [3]Node
	Slaves  [3]Node
}
