/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/
package cluster

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math/rand/v2"
	"os"
	"time"
	"topiik/internal/config"
	"topiik/internal/util"
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
	if nodeStatus.Role == RAFT_LEADER {
		nodeStatus.Role = RAFT_FOLLOWER
	}

	// update Raft Heartbeat
	nodeStatus.Heartbeat = uint16(rand.IntN(int(serverConfig.RaftHeartbeatMax-serverConfig.RaftHeartbeatMin))) + serverConfig.RaftHeartbeatMin
	nodeStatus.HeartbeatAt = time.Now().UnixMilli()

	var entryType int8 // tow bytes of command
	if len(entry) >= 1 {
		entryTypeByte := entry[:1]
		byteBuf := bytes.NewBuffer(entryTypeByte)
		err := binary.Read(byteBuf, binary.LittleEndian, &entryType)
		if err != nil {
			tLog.Err(err)
		}

		if entryType == ENTRY_TYPE_DEFAULT { // append controller address
			// log.Info().Msgf("appendEntry() Leader addr:%s", string(entry[1:]))
			nodeStatus.LeaderControllerAddr = string(entry[1:])
		} else if entryType == ENTRY_TYPE_METADATA { // append cluster metadata
			tLog.Info().Msg("appendEntry: metadata")
			var clusterData = &Cluster{}
			err := json.Unmarshal(entry[1:], clusterData) // verify
			if err != nil {
				tLog.Err(err)
				return err
			}
			clusterPath := GetClusterFilePath()
			exist, _ := util.PathExists(clusterPath)
			if exist {
				err = os.Truncate(clusterPath, 0) // TODO: backup first
				if err != nil {
					tLog.Err(err)
					return err
				}
			}
			err = os.WriteFile(clusterPath, entry[1:], 0644)
			if err != nil {
				tLog.Err(err)
				return err
			}
		} else if entryType == ENTRY_TYPE_PARTITION {
			tLog.Info().Msg("appendEntry: partittion")
			var partitions = make(map[string]Partition)
			err := json.Unmarshal(entry[1:], &partitions) // verify
			if err != nil {
				tLog.Err(err)
				return err
			}
			filePath := GetPatitionFilePath()
			exist, _ := util.PathExists(filePath)
			if exist {
				os.Remove(filePath) //TODO: move to archive
			}
			err = util.WriteBinaryFile(filePath, entry[1:])
			if err != nil {
				tLog.Err(err)
				return err
			}
		}
	}

	return nil
}
