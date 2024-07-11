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
)

func appendEntry(pieces []string, serverConfig *config.ServerConfig) error {
	//fmt.Println(pieces)
	if len(pieces) == 2 {
		if pieces[0] == "CONTROLLER" {
			fmt.Println("append_entry::appendEntry() CONTROLLER")
			var tmpMap map[string]NodeSlim
			err := json.Unmarshal([]byte(pieces[1]), &tmpMap)
			if err != nil {
				return err
			}
			fmt.Println(tmpMap)
			controllerPath := GetControllerFilePath()
			err = os.WriteFile(controllerPath, []byte(pieces[1]), 0644)
			if err != nil {
				return err
			}
			return nil
		} else if pieces[0] == "WORKER" {
			fmt.Println("append_entry::appendEntry() WORKER")
		} else if pieces[0] == "PARTITION" {
			fmt.Println("append_entry::appendEntry() PARTITION")
		}
	}
	nodeStatus.Heartbeat = uint16(rand.IntN(int(serverConfig.RaftHeartbeatMax-serverConfig.RaftHeartbeatMin))) + serverConfig.RaftHeartbeatMin
	nodeStatus.HeartbeatAt = time.Now().UnixMilli()
	return nil
}