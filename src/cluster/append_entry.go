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
	"topiik/logger"
)

/***
**
**
** Parameter: 1 byte of entry type + entry
**
**
***/
func appendEntry(entry []byte, serverConfig *config.ServerConfig) error {
	var log = logger.Get()
	var entryType int8 // tow bytes of command
	if len(entry) >= 1 {
		entryTypeByte := entry[:1]
		byteBuf := bytes.NewBuffer(entryTypeByte)
		err := binary.Read(byteBuf, binary.LittleEndian, &entryType)
		if err != nil {
			log.Err(err)
		}

		if entryType == ENTRY_TYPE_DEFAULT { // append controller address
			// log.Info().Msgf("appendEntry() Leader addr:%s", string(entry[1:]))
			nodeStatus.LeaderControllerAddr = string(entry[1:])
		} else if entryType == ENTRY_TYPE_METADATA { // append cluster metadata
			log.Info().Msg("appendEntry() metadata")
			var clusterData = &Cluster{}
			err := json.Unmarshal(entry[1:], clusterData) // verify
			if err != nil {
				log.Err(err)
				return err
			}
			clusterPath := GetClusterFilePath()
			exist, _ := util.PathExists(clusterPath)
			if exist {
				err = os.Truncate(clusterPath, 0) // TODO: backup first
				if err != nil {
					log.Err(err)
					return err
				}
			}
			err = os.WriteFile(clusterPath, entry[1:], 0644)
			if err != nil {
				log.Err(err)
				return err
			}
		}
	}

	// update Raft Heartbeat
	nodeStatus.Heartbeat = uint16(rand.IntN(int(serverConfig.RaftHeartbeatMax-serverConfig.RaftHeartbeatMin))) + serverConfig.RaftHeartbeatMin
	nodeStatus.HeartbeatAt = time.Now().UnixMilli()
	return nil
}
