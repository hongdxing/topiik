/***
**
**/

package cluster

const (
	RAFT_LEADER     = 1
	RAFT_FOLLOWER   = 2
	RAFT_CANDIDATOR = 3
)

const (
	ROLE_CONTROLLER = "CONTROLLER"
	ROLE_WORKER     = "WORKER"
)

const (
	RPC_VOTE         int16 = 1
	RPC_APPENDENTRY  int16 = 2
	RPC_ADD_NODE     int16 = 3
	CLUSTER_JOIN_ACK int16 = 4
)

const (
	ENTRY_TYPE_DEFAULT  int8 = 1
	ENTRY_TYPE_METADATA int8 = 2
)

const (
	VOTE_ACCEPTED = "A"
	VOTE_REJECTED = "R"
)

const SLOTS = 1024

const CONTROLLER_FORWORD_PORT = "9302"
