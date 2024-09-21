/*
* ©2024 www.topiik.com
* author: Duan Hongxing
* data: 7 Jul 2024
*
 */

package server

import (
	"fmt"
	"time"
	"topiik/cluster"
	"topiik/internal/consts"
	"topiik/node"
)

/*
*
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

	if node.IsWorker() {
		if cluster.GetNodeStatus().RaftRole == cluster.RAFT_LEADER {
			return consts.VOTE_REJECTED + ":L"
		}
		if cluster.GetNodeStatus().RaftRole != cluster.RAFT_FOLLOWER {
			return consts.VOTE_REJECTED + ":F"
		}
		if cluster.GetTerm() > cTerm {
			return consts.VOTE_REJECTED + ":L" // if current node version greater than candidate's version, reject as Leader reject
		}
	} else {
		fmt.Println("worker vote")
		// if the woker still can sense the controller Leader, then reject
		if time.Now().UTC().UnixMilli() < cluster.GetNodeStatus().HeartbeatAt+int64(cluster.GetNodeStatus().Heartbeat) {
			return consts.VOTE_REJECTED + ":W"
		}
		return consts.VOTE_ACCEPTED + ":W"
	}

	// (lastTermV > lastTermC) || ((lastTermV == lastTermC) && (lastIndexV > lastIndexC))
	if cluster.GetTerm() > cTerm {
		return consts.VOTE_REJECTED + ":T"
	} else {
		return consts.VOTE_ACCEPTED + ":A"
	}
}
