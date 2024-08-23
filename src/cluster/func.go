package cluster

func SetRole(role uint8) {
	nodeStatus.Role = role
}

func SetLeaderCtlAddr(addr string) {
	nodeStatus.LeaderControllerAddr = addr
}

func SetPtnInfo(ptnInfo *PartitionInfo) {
	partitionInfo = ptnInfo
}

func SetHeartbeat(heartbeat uint16, heartbeatAt int64) {
	nodeStatus.Heartbeat = heartbeat
	nodeStatus.HeartbeatAt = heartbeatAt
}
