/*
* @author: Duan Hongxing
* @date: 23 Aug, 2024
* @desc:
*	Cluster info implementation
 */

package cluster

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"topiik/node"
)

/*
* Client issue command ADD-WORKER or ADD-CONTROLLER
* After the target node accepted, Controller add node to metadata
*
 */
func AddNode(ndId string, addr string, addr2 string, role string) (err error) {
	if strings.ToUpper(role) == ROLE_CONTROLLER {
		clusterInfo.Ctls[ndId] = node.NodeSlim{Id: ndId, Addr: addr, Addr2: addr2}
	} else {
		worker := node.NodeSlim{Id: ndId, Addr: addr, Addr2: addr2}
		clusterInfo.Wkrs[ndId] = worker
	}

	/* save cluster to disk */
	clusterPath := GetClusterFilePath()
	buf, err := json.Marshal(clusterInfo)
	if err != nil {
		return errors.New("update cluster failed")
	}
	err = os.Truncate(clusterPath, 0) // TODO: myabe need backup first
	if err != nil {
		return errors.New("update cluster failed")
	}
	err = os.WriteFile(clusterPath, buf, 0664) // save back controller file
	if err != nil {
		return errors.New("update cluster failed")
	}

	/* cluster meta changed, pending sync to follower(s) */
	UpdatePendingAppend()

	return nil
}
