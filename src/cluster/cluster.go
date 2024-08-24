/*
* author: Duan Hongxing
* date: 23 Aug, 2024
* desc:
*	Cluster info implementation
 */

package cluster

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"topiik/internal/util"
	"topiik/node"
)

/*
* Client issue command ADD-WORKER or ADD-CONTROLLER
* After the target node accepted, Controller add node to metadata
*
 */
func AddNode(ndId string, addr string, addr2 string, role string) (err error) {
	if strings.ToUpper(role) == node.ROLE_CONTROLLER {
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
	notifyMetadataChanged()
	notifyPtnChanged()

	return nil
}

func SetClusterInfo(cluster *Cluster) {
	clusterInfo = cluster
}

func SaveClusterInfo(data []byte) (err error) {
	fpath := GetClusterFilePath()
	exist, _ := util.PathExists(fpath)
	if exist {
		err = os.Truncate(fpath, 0) // TODO: backup first
		if err != nil {
			l.Err(err)
			return err
		}
	}
	err = os.WriteFile(fpath, data, 0644)
	if err != nil {
		l.Err(err)
		return err
	}
	return nil
}
