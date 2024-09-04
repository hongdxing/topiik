// author: Duan Hongxing
// date: 23 Aug, 2024
// desc:	Cluster info implementation

package cluster

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

// Client issue command ADD-WORKER or ADD-CONTROLLER
// After the target node accepted, Controller add node to metadata
func AddNode(ndId string, addr string, addr2 string, role string, ptnId string) (err error) {
	if strings.ToUpper(role) == node.ROLE_CONTROLLER {
		controllerInfo.Nodes[ndId] = node.NodeSlim{Id: ndId, Addr: addr, Addr2: addr2}
		saveControllerInfo()
	} else if strings.ToUpper(role) == node.ROLE_WORKER {
		workerInfo.Nodes[ndId] = node.NodeSlim{Id: ndId, Addr: addr, Addr2: addr2}
		saveWorkerInfo()
		addNode2Partition(ptnId, ndId)
	} else {
		return errors.New("")
	}

	notifyControllerChanged()
	notifyWorkerChanged()
	notifyPtnChanged()

	return nil
}

// Remove node, Controller or Worker from cluster
// Syntax: REMOVE-NODE nodeId
func RemoveNode(ndId string) (err error) {
	if nd, ok := controllerInfo.Nodes[ndId]; ok {
		/* if there is only one Controller in cluster, then reject */
		if len(controllerInfo.Nodes) == 1 {
			return errors.New(resp.RES_REJECTED)
		}
		/* if trying to remove current node and current node is Controller Leader, then reject */
		if nd.Id == node.GetNodeInfo().Id && nodeStatus.Role == RAFT_LEADER {
			return errors.New(resp.RES_REJECTED)
		}
		delete(controllerInfo.Nodes, ndId)
		saveControllerInfo()
		notifyControllerChanged()
	} else if _, ok := workerInfo.Nodes[ndId]; ok {
		/* if this is the only worker node, then reject */
		if len(workerInfo.Nodes) == 1 {
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
		delete(workerInfo.Nodes, ndId)

		go rpcRemoveNode(ndId)

		saveWorkerInfo()
		notifyWorkerChanged()

		err = savePartition()
		notifyPtnChanged()
		if err != nil {
			return err
		}
	} else {
		return errors.New(resp.RES_NIL)
	}
	return nil
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

/*
func saveClusterInfo() (err error) {
	data, err := json.Marshal(clusterInfo)
	if err != nil {
		l.Err(err).Msgf("cluster::saveClusterInfo %s", err.Error())
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
*/

// RPC to remove cluster info of the node
func rpcRemoveNode(ndId string) {
	var addr2 string
	if nd, ok := controllerInfo.Nodes[ndId]; ok {
		addr2 = nd.Addr2
	} else if nd, ok := workerInfo.Nodes[ndId]; ok {
		addr2 = nd.Addr2
	}
	if addr2 == "" {
		return
	}

	var buf []byte
	var bbuf = new(bytes.Buffer) // int to byte buf
	_ = binary.Write(bbuf, binary.LittleEndian, consts.RPC_REMOVE_NODE)
	buf = append(buf, bbuf.Bytes()...)
	buf = append(buf, []byte(ndId)...)

	// Enocde
	buf, err := proto.EncodeB(buf)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return
	}

	// Send
	conn, err := util.PreapareSocketClient(addr2)
	if err != nil {
		l.Warn().Msgf("cluster::rpcRemoveNode Cannot connect to the addr2 %s", addr2)
		return
	}
	_, err = conn.Write(buf)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return
	}
}
