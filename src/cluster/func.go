package cluster

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"topiik/internal/consts"
	"topiik/internal/util"
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
		worker := node.NodeSlim{Id: nodeId, Addr: addr, Addr2: addr2}
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

/*
* Client issue command NEW-PARTITION
*
 */
func NewPartition(ptnCount int) (ptnIds []string, err error) {
	if len(partitionInfo.ClusterId) == 0 { // brand new cluster without partition yet
		for i := 0; i < int(ptnCount); i++ {
			var from int
			var to int
			from = i * (consts.SLOTS / ptnCount) // p=2--> i=0: 0, i=1: 512

			if i == (ptnCount - 1) {
				to = consts.SLOTS - 1
			} else {
				to = (i+1)*(consts.SLOTS/ptnCount) - 1 // p=2--> i=0: 511, i=1: 1024
			}
			slot := node.Slot{From: uint16(from), To: uint16(to)}

			ptnId := util.RandStringRunes(10)
			ptnIds = append(ptnIds, ptnId)
			partitionInfo.PtnMap[ptnId] = &node.Partition{
				Id:      ptnId,
				NodeSet: make(map[string]*node.NodeSlim),
				Slots:   []node.Slot{slot},
			}
		}
	} else { // having existing partition(s), TODO
		err = errors.New("cannot create new partition")
		return ptnIds, err
	}

	// persist
	filePath := GetPatitionFilePath()
	exist, err := util.PathExists(filePath)
	if err != nil {
		l.Err(err).Msgf("scale: %s", err.Error())
		return ptnIds, err
	}
	if exist { // rename to old for backup
		err = os.Rename(filePath, filePath+"old")
		if err != nil {
			l.Err(err).Msgf("scale: %s", err.Error())
			return ptnIds, err
		}
	}
	data, err := json.Marshal(partitionInfo)
	if err != nil {
		l.Err(err).Msgf("scale: %s", err.Error())
		return ptnIds, err
	}

	err = util.WriteBinaryFile(filePath, data)
	if err != nil {
		l.Err(err).Msgf("scale: %s", err.Error())
		return ptnIds, err
	}

	return ptnIds, err
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
