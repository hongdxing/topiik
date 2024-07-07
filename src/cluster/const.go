/***
**
**/

package cluster

const (
	RAFT_LEADER     = 1
	RAFT_CANDIDATOR = 2
	RAFT_FOLLOWER   = 3
)

const (
	ROLE_CONTROLLER = "CONTROLLER"
	ROLE_WORKER     = "WORKER"
)

const (
	RPC_APPENDENTRY = "APPENDENTRY"
)

const (
	VOTE_ACCEPTED = "A"
	VOTE_REJECTED = "R"
)
