/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/
package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math/rand/v2"
	"os"
	"time"
	"topiik/cluster"
	"topiik/internal/config"
	"topiik/internal/util"
	"topiik/node"
)

/*
* Raft Append Entry
*
* Parameter:
* 	entry: 1 byte of entry type + entry
*
*
 */
func appendEntry(entry []byte, serverConfig *config.ServerConfig) error {
	// In case of multi Leader, if node can receive appendEntry,
	// and role is RAFT_LEADER, then step back
	if cluster.IsNodeController() && cluster.GetNodeStatus().Role == cluster.RAFT_LEADER {
		//cluster.GetNodeStatus().Role = cluster.RAFT_FOLLOWER
		cluster.SetRole(cluster.RAFT_FOLLOWER)
		go cluster.RequestVote()
	}

	// update Raft Heartbeat
	//nodeStatus.Heartbeat = uint16(rand.IntN(int(serverConfig.RaftHeartbeatMax-serverConfig.RaftHeartbeatMin))) + serverConfig.RaftHeartbeatMin
	//nodeStatus.HeartbeatAt = time.Now().UnixMilli()
	heartbeat := uint16(rand.IntN(int(serverConfig.RaftHeartbeatMax-serverConfig.RaftHeartbeatMin))) + serverConfig.RaftHeartbeatMin
	heartbeatAt := time.Now().UnixMilli()
	cluster.SetHeartbeat(heartbeat, heartbeatAt)

	var entryType int8 // one byte of command
	if len(entry) >= 1 {
		entryTypeByte := entry[:1]
		byteBuf := bytes.NewBuffer(entryTypeByte)
		err := binary.Read(byteBuf, binary.LittleEndian, &entryType)
		if err != nil {
			l.Err(err)
		}

		if entryType == cluster.ENTRY_TYPE_DEFAULT { // append controller address
			// log.Info().Msgf("appendEntry() Leader addr:%s", string(entry[1:]))
			//nodeStatus.LeaderControllerAddr = string(entry[1:])
			cluster.SetLeaderCtlAddr(string(entry[1:]))
		} else if entryType == cluster.ENTRY_TYPE_PTN { // append worker followers
			node.SetPtn(entry[1:])
		} else if entryType == cluster.ENTRY_TYPE_METADATA { // append cluster metadata
			l.Info().Msg("rpc_append_entry::appendEntry metadata begin")
			var clusterData = &cluster.Cluster{}
			err := json.Unmarshal(entry[1:], clusterData) // verify
			if err != nil {
				l.Err(err)
				return err
			}
			clusterPath := cluster.GetClusterFilePath()
			exist, _ := util.PathExists(clusterPath)
			if exist {
				err = os.Truncate(clusterPath, 0) // TODO: backup first
				if err != nil {
					l.Err(err)
					return err
				}
			}
			err = os.WriteFile(clusterPath, entry[1:], 0644)
			if err != nil {
				l.Err(err)
				return err
			}
			l.Info().Msg("rpc_append_entry::appendEntry metadata end")
		} else if entryType == cluster.ENTRY_TYPE_PTNS {
			l.Info().Msg("rpc_append_entry::appendEntry partittion begin")
			var ptnInfo cluster.PartitionInfo
			err := json.Unmarshal(entry[1:], &ptnInfo) // verify
			if err != nil {
				l.Err(err).Msg(err.Error())
				return err
			}
			//partitionInfo = &ptnInfo
			cluster.SetPtnInfo(&ptnInfo)
			filePath := cluster.GetPatitionFilePath()
			exist, _ := util.PathExists(filePath)
			if exist {
				os.Rename(filePath, filePath+"old") // rename
			}
			err = util.WriteBinaryFile(filePath, entry[1:])
			if err != nil {
				l.Err(err).Msg(err.Error())
				return err
			}
			l.Info().Msg("rpc_append_entry::appendEntry partittion end")
		}
	}

	return nil
}
