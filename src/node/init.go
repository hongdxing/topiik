package node

import (
	"topiik/internal/config"
	"topiik/internal/logger"
)

var nodeInfo *Node

var partition Partition

var serverConfig *config.ServerConfig

var l = logger.Get()
