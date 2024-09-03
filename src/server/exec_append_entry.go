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
	if node.IsController() && cluster.GetNodeStatus().Role == cluster.RAFT_LEADER {
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
			//l.Info().Msgf("appendEntry() Leader addr:%s", string(entry[1:]))
			//nodeStatus.LeaderControllerAddr = string(entry[1:])
			cluster.SetLeaderCtlAddr(string(entry[1:]))
		} else if entryType == cluster.ENTRY_TYPE_PTN { // append worker followers
			node.SetPtn(entry[1:])
		} else if entryType == cluster.ENTRY_TYPE_METADATA { // append cluster metadata
			l.Info().Msg("rpc_append_entry::appendEntry metadata begin")
			var clusterInfo = &cluster.Cluster{}
			err := json.Unmarshal(entry[1:], clusterInfo) // verify
			if err != nil {
				l.Err(err)
				return err
			}
			/* set cluster info in memory */
			//cluster.SetClusterInfo(clusterInfo)
			l.Info().Msg("rpc_append_entry::appendEntry metadata end")
		} else if entryType == cluster.ENTRY_TYPE_CTL {
			l.Info().Msg("rpc_append_entry::appendEntry controller begin")
			var controllerInfo = &cluster.NodesInfo{}
			err := json.Unmarshal(entry[1:], controllerInfo)
			if err != nil {
				l.Err(err)
				return err
			}
			cluster.SetControllerInfo(controllerInfo)
			l.Info().Msg("rpc_append_entry::appendEntry controller end")
		} else if entryType == cluster.ENTRY_TYPE_WRK {
			l.Info().Msg("rpc_append_entry::appendEntry worker begin")
			var workerInfo = &cluster.NodesInfo{}
			err := json.Unmarshal(entry[1:], workerInfo)
			if err != nil {
				l.Err(err)
				return err
			}
			cluster.SetWorkerInfo(workerInfo)
			l.Info().Msg("rpc_append_entry::appendEntry worker end")
		} else if entryType == cluster.ENTRY_TYPE_PTNS {
			l.Info().Msg("rpc_append_entry::appendEntry partition begin")
			var ptnInfo cluster.PartitionInfo
			err := json.Unmarshal(entry[1:], &ptnInfo) // verify
			if err != nil {
				l.Err(err).Msg(err.Error())
				return err
			}
			//partitionInfo = &ptnInfo
			cluster.SetPtnInfo(&ptnInfo)
			fpath := cluster.GetPatitionFilePath()
			exist, _ := util.PathExists(fpath)
			if exist {
				os.Rename(fpath, fpath+"old") // rename
			}
			err = util.WriteBinaryFile(fpath, entry[1:])
			if err != nil {
				l.Err(err).Msg(err.Error())
				return err
			}
			l.Info().Msg("rpc_append_entry::appendEntry partition end")
		}
	}

	return nil
}
