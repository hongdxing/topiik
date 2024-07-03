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

/*** Cluster init sub commands ***/
const(
	CLUSTER_INIT_PRE_CHECK = "__CLUSTER_INIT_PRE_CHECK__"
	CLUSTER_INIT_CONFIRM = "__CLUSTER_INIT_CONFIRM__"
)

/*** Cluster init result ***/
const (
	RES_CLUSTER_INIT_OK            = "OK"
	RES_CLUSTER_INIT_FAILED        = "FAILED"
	RES_CLUSTER_INIT_NETWORK_ISSUE = "NETWORK"
)
