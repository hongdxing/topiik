/***
** author: duan hongxing
** date: 21 Jul 2024
** desc:
**
**/

package server

import (
	"errors"
	"strings"
	"topiik/cluster"
	"topiik/node"
	"topiik/resp"
)

/*
* Controller send ADD-NODE RPC to current node
* Parameters:
*	- pieces[0]: clusterId
*	- pieces[1]: role
 */
func addNode(pieces []string) (string, error) {
	if len(pieces) != 2 {
		return "", errors.New(resp.RES_SYNTAX_ERROR)
	}
	clusterId := pieces[0]
	role := pieces[1]

	node.JoinCluster(clusterId)

	// if join controller succeed, will start to RequestVote
	if strings.ToUpper(role) == cluster.ROLE_CONTROLLER {
		go cluster.RequestVote()
	}
	// return nodeId to controller
	return node.GetNodeInfo().Id, nil
}
