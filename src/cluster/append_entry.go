/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/
package cluster

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
