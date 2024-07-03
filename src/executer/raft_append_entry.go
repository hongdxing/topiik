/***
* author: duan hongxing
* date: 22 Jun 2024
* desc:
* 		When received appendEntry, this node need to do:
*		1): update nodeStatus.Heartbeat with random number(200, 500) milli seconds
*		2): append log(data)
***/

package executer

import (
	"math/rand/v2"
	"time"
	"topiik/internal/config"
	"topiik/raft"
)

func appendEntry(serverConfig *config.ServerConfig, nodeStatus *raft.NodeStatus) {
	nodeStatus.Heartbeat = uint16(rand.IntN(int(serverConfig.RaftHeartbeatMax-serverConfig.RaftHeartbeatMin))) + serverConfig.RaftHeartbeatMin
	nodeStatus.HeartbeatAt = time.Now().UnixMilli()
}
