/***
** author: duan hongxing
** date: 21 Jul 2024
** desc:
**
**/

package cluster

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

/*
** Parameters:
**	- pieces[0]: clusterId
**	- pieces[1]: role
**
**
 */
func addNode(pieces []string) (string, error) {
	if len(pieces) != 2 {
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	clusterId := pieces[0]
	role := pieces[1]

	// update node cluster id
	nodeInfo.ClusterId = clusterId

	nodePath := GetNodeFilePath()
	buf, err := json.Marshal(nodeInfo)
	if err != nil {
		return "", errors.New("update node failed")
	}
	err = os.Truncate(nodePath, 0) // TODO: myabe need backup first
	if err != nil {
		return "", errors.New("update node failed")
	}
	err = os.WriteFile(nodePath, buf, 0664) // save back controller file
	if err != nil {
		return "", errors.New("update node failed")
	}

	// if join controller succeed, will start to RequestVote
	if strings.ToUpper(role) == ROLE_CONTROLLER {
		go RequestVote()
	}
	// return nodeId to controller
	return nodeInfo.Id, nil
}
