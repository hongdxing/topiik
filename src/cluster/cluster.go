// author: Duan Hongxing
// date: 23 Aug, 2024
// desc:	Cluster info implementation

package cluster

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
)

func SetRole(role uint8) {
	nodeStatus.RaftRole = role
}

func SetLeaderCtlAddr(addr string) {
	nodeStatus.LeaderControllerAddr = addr
}

func SetWorkerGroupInfo(data []byte) error {
	err := json.Unmarshal(data, &partitionInfo)
	fmt.Println(partitionInfo)
	savePartitions()
	return err
}

func SetHeartbeat(heartbeat uint16, heartbeatAt int64) {
	nodeStatus.Heartbeat = heartbeat
	nodeStatus.HeartbeatAt = heartbeatAt
}

func GetPtnLeaders() (workers []node.NodeSlim) {
	for _, ptn := range partitionInfo.Ptns {
		leader := ptn.Nodes[ptn.LeaderNodeId]
		workers = append(workers, leader)
	}
	return workers
}

func GetPtnLeader(ndId string) (leader node.NodeSlim) {
	for _, ptn := range partitionInfo.Ptns {
		if _, ok := ptn.Nodes[ndId]; ok {
			if ptn.LeaderNodeId != "" {
				leader = ptn.Nodes[ptn.LeaderNodeId]
				break
			}
		}
	}
	// safe? what if leader still empty???
	return leader
}

func GetPtnByNodeId(ndId string) (partition Partition) {
	for _, ptn := range partitionInfo.Ptns {
		if _, ok := ptn.Nodes[ndId]; ok {
			partition = *ptn
			break
		}
	}
	return partition
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
	//if nd, ok := controllerInfo.Nodes[ndId]; ok {
	//	addr2 = nd.Addr2
	//}
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
