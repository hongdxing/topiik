/***
** author: duan hongxing
** data: 7 Jul 2024
** desc:
**
**/

package cluster

import (
	"fmt"
	"time"
)

/*
**
  - Return: [A|R]:[A|L|F|T]
  - return value seperated by :
  - 1) first part, Aaccept or Rejected
  - 2) sencond part is reason
    -- if first part is A(ccespted), then sencond part also A, actually meanningless
    -- if first part is R(ejected), then second part is:
    ---- 1) L(eader), mean there is a Leader already
    ---- 2) F(ollower), mean current node is Follower, Follower not allow to vote
    ---- 3) T(erm), means Term of current node is bigger than request node
*/
func vote(cTerm int) string {
	//fmt.Printf("vote():: current node role: %v\n", nodeStatus.Role)

	if IsNodeController() {
		if nodeStatus.Role == RAFT_LEADER {
			return VOTE_REJECTED + ":L"
		}
		if nodeStatus.Role != RAFT_FOLLOWER {
			return VOTE_REJECTED + ":F"
		}
		if clusterInfo.Ver > uint(cTerm) {
			return VOTE_REJECTED + ":L" // if current node version greater than candidate's version, reject as Leader reject
		}
	} else {
		fmt.Println("worker vote")
		// if the woker still can sense the controller Leader, then reject
		if time.Now().UTC().UnixMilli() < nodeStatus.HeartbeatAt+int64(nodeStatus.Heartbeat) {
			return VOTE_REJECTED + ":W"
		}
		return VOTE_ACCEPTED + ":W"
	}

	// (lastTermV > lastTermC) || ((lastTermV == lastTermC) && (lastIndexV > lastIndexC))
	if nodeStatus.Term > uint(cTerm) {
		return VOTE_REJECTED + ":T"
	} else {
		return VOTE_ACCEPTED + ":A"
	}
}
