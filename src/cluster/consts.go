/***
** author: duan hongxing
** date: 30 Jun 2024
** desc:
**	return keys
**
**/

package cluster

/*** Node Role ***/
const (
	NODE_ROLE_MASTER     string = "M"
	NODE_ROLE_FOLLOWER   string = "F"
	NODE_ROLE_STANDALONE string = "S" // Not joined cluster
)

/*** Response message ***/
const (
	RES_INVALID_NUM_OF_NODE string = "INVALID_NUM_OF_NODE"
)

/*** Cluster init result ***/
const (
	CLUSTER_INIT_OK            = "OK"
	CLUSTER_INIT_FAILED        = "FAILED"
	CLUSTER_INIT_NETWORK_ISSUE = "NETWORK"
)
