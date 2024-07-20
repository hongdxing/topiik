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
	CLUSTER_JOIN_ACK int16 = 3
)


const (
	ENTRY_TYPE_DEFAULT  int8 = 1
	ENTRY_TYPE_METADATA int8 = 2
)

const (
	VOTE_ACCEPTED = "A"
	VOTE_REJECTED = "R"
)

const CONTROLLER_FORWORD_PORT = "9302"
