/***
** author: duan hongxing
** date: 21 Jul 2024
** desc:
**
**/

package cluster

import (
	"errors"
	"strings"
	"topiik/node"
)

/*
* Controller send ADD-NODE RPC to current node
* Parameters:
*	- pieces[0]: clusterId
*	- pieces[1]: role
 */
func addNode(pieces []string) (string, error) {
	if len(pieces) != 2 {
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	clusterId := pieces[0]
	role := pieces[1]

	node.JoinCluster(clusterId)

	// if join controller succeed, will start to RequestVote
	if strings.ToUpper(role) == ROLE_CONTROLLER {
		go RequestVote()
	}
	// return nodeId to controller
	return node.GetNodeInfo().Id, nil
}
