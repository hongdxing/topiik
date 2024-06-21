package executor

import (
	"fmt"
	"topiik/raft"
)

func vote(cTerm int, nodeStatus *raft.NodeStatus) string {
	fmt.Printf("vote():: current node role: %v\n", nodeStatus.Role)
	if nodeStatus.Role != raft.ROLE_CANDIDATOR {
		return RES_REJECTED
	}
	// (lastTermV > lastTermC) || ((lastTermV == lastTermC) && (lastIndexV > lastIndexC))
	if nodeStatus.Term > uint(cTerm) {
		return RES_REJECTED
	} else {
		return RES_ACCEPTED
	}
}
