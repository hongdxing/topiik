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
	"topiik/resp"
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

/*
* Remove node, Controller or Worker from cluster
* Syntax: REMOVE-NODE nodeId
 */
func RemoveNode(ndId string) (err error) {
	if nd, ok := clusterInfo.Ctls[ndId]; ok {
		/* if there is only one Controller in cluster, then reject */
		if len(clusterInfo.Ctls) == 1 {
			return errors.New(resp.RES_REJECTED)
		}
		/* if trying to remove current node and current node is Controller Leader, then reject */
		if nd.Id == node.GetNodeInfo().Id && nodeStatus.Role == RAFT_LEADER {
			return errors.New(resp.RES_REJECTED)
		}
		delete(clusterInfo.Ctls, ndId)
		notifyMetadataChanged()
		err = saveClusterInfo()
		if err != nil {
			return err
		}
	} else if _, ok := clusterInfo.Wkrs[ndId]; ok {
		/* if this is the only worker node, then reject */
		if len(clusterInfo.Wkrs) == 1 {
			return errors.New(resp.RES_REJECTED)
		}
		/* if is partition leader, then reject */
		for _, ptn := range partitionInfo.PtnMap {
			if ndId == ptn.LeaderNodeId {
				return errors.New(resp.RES_REJECTED)
			}
		}

		for _, ptn := range partitionInfo.PtnMap {
			if _, ok := ptn.NodeSet[ndId]; ok {
				if len(ptn.NodeSet) == 1 {
					return errors.New(resp.RES_REJECTED)
				}
			}
			delete(ptn.NodeSet, ndId)
		}
		delete(clusterInfo.Wkrs, ndId)
		notifyMetadataChanged()

		err = saveClusterInfo()
		if err != nil {
			return err
		}

		err = savePartition()
		if err != nil {
			return err
		}
	} else {
		return errors.New(resp.RES_NIL)
	}
	return nil
}

func SetClusterInfo(cluster *Cluster) {
	clusterInfo = cluster
	saveClusterInfo()
}

func SetRole(role uint8) {
	nodeStatus.Role = role
}

func SetLeaderCtlAddr(addr string) {
	nodeStatus.LeaderControllerAddr = addr
}

func SetPtnInfo(ptnInfo *PartitionInfo) {
	partitionInfo = ptnInfo
}

func SetHeartbeat(heartbeat uint16, heartbeatAt int64) {
	nodeStatus.Heartbeat = heartbeat
	nodeStatus.HeartbeatAt = heartbeatAt
}

/*pivate func----------------------------------------------------------------*/

func saveClusterInfo() (err error) {
	data, err := json.Marshal(clusterInfo)
	if err != nil {
		l.Err(err).Msgf("cluster::RemoveNode %s", err.Error())
		return err
	}

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
