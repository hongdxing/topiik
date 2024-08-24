/*
* author: duan hongxing
* date: 21 Jul 2024
* desc:
*
 */

package server

import (
	"errors"
	"strings"
	"topiik/cluster"
	"topiik/node"
	"topiik/resp"
)

/*
* Controller send ADD-WORKER|ADD-CONTROLLER RPC to current node
* Parameters:
*	- pieces[0]: clusterId
*	- pieces[1]: role
 */
func addNode(pieces []string) (string, error) {
	/* validate: make sure node not belongs to any cluster yet */
	if node.GetNodeInfo().ClusterId != "" {
		return "", errors.New("target node already in cluster: " + pieces[0])
	}

	if len(pieces) != 2 {
		return "", errors.New(resp.RES_SYNTAX_ERROR)
	}
	clusterId := pieces[0]
	role := strings.ToUpper(pieces[1])

	if role != node.ROLE_CONTROLLER && role != node.ROLE_WORKER {
		return "", errors.New("invalid role: " + role)
	}

	node.JoinCluster(clusterId, role)

	/* if join controller succeed, will start to RequestVote */
	if role == node.ROLE_CONTROLLER {
		go cluster.RequestVote()
	}
	// return nodeId to controller
	return node.GetNodeInfo().Id, nil
}
