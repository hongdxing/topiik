/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/
package cluster

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
	"topiik/internal/config"
	"topiik/internal/util"
)

/***
**
**
**
**
**
***/
func appendEntry(pieces []string, serverConfig *config.ServerConfig) error {
	if len(pieces) == 2 {
		/*
			var tmpMap map[string]NodeSlim
			err := json.Unmarshal([]byte(pieces[1]), &tmpMap)
			if err != nil {
				return err
			}
			fmt.Println(tmpMap)
		*/

		if pieces[0] == "METADATA" {
			fmt.Println("append_entry::appendEntry() cluster metadata")
			err := json.Unmarshal([]byte(pieces[1]), clusterInfo) // verify
			if err != nil {
				return err
			}
			clusterPath := GetClusterFilePath()
			exist, _ := util.PathExists(clusterPath)
			if exist {
				err = os.Truncate(clusterPath, 0) // TODO: backup first
				if err != nil {
					return err
				}
			}
			err = os.WriteFile(clusterPath, []byte(pieces[1]), 0644)
			if err != nil {
				return err
			}
		}
	}
	// update Raft Heartbeat
	nodeStatus.Heartbeat = uint16(rand.IntN(int(serverConfig.RaftHeartbeatMax-serverConfig.RaftHeartbeatMin))) + serverConfig.RaftHeartbeatMin
	nodeStatus.HeartbeatAt = time.Now().UnixMilli()
	return nil
}
