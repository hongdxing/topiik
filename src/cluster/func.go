package cluster

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"topiik/node"
)

/*
* Client issue command ADD-NODE
* After the target node accepted, Controller add node to metadata
*
 */
func AddNode(nodeId string, addr string, addr2 string, role string) (err error) {
	if strings.ToUpper(role) == ROLE_CONTROLLER {
		clusterInfo.Ctls[nodeId] = node.NodeSlim{Id: nodeId, Addr: addr, Addr2: addr2}
	} else {
		worker := Worker{Id: nodeId, Addr: addr, Addr2: addr2}
		clusterInfo.Wkrs[nodeId] = worker
	}
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
	// cluster meta changed, pending to sync to follower(s)
	UpdatePendingAppend()

	return nil
}

func GetPartitionInfo() PartitionInfo {
	return *partitionInfo
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
