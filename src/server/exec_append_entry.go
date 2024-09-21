// author: duan hongxing
// data: 3 Jul 2024
package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math/rand/v2"
	"time"
	"topiik/cluster"
	"topiik/internal/config"
	"topiik/node"
)

// Raft Append Entry
// Parameter:
//
//	entry: 1 byte of entry type + entry
func appendEntry(entry []byte, serverConfig *config.ServerConfig) error {
	// In case of multi Leader, if node can receive appendEntry,
	// and role is RAFT_LEADER, then step back
	if node.IsWorker() && cluster.GetNodeStatus().RaftRole == cluster.RAFT_LEADER {
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
		} else if entryType == cluster.ENTRY_TYPE_WRKGRP {
			l.Info().Msg("rpc_append_entry::appendEntry worker group begin")
			err := cluster.SetWorkerGroupInfo(entry[1:])
			l.Err(err)
			l.Info().Msg("rpc_append_entry::appendEntry worker group end")
		}
	}

	return nil
}
