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
	ENTRY_TYPE_DEFAULT  int8 = 1
	ENTRY_TYPE_METADATA int8 = 2
	ENTRY_TYPE_PTN      int8 = 3
	ENTRY_TYPE_PTNS     int8 = 4
	ENTRY_TYPE_PTN_LDR  int8 = 5
)

const (
	VOTE_ACCEPTED = "A"
	VOTE_REJECTED = "R"
)

const CONTROLLER_FORWORD_PORT = "9302"
