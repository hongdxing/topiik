/*
* author: duan hongxing
* date: 6 July, 2024
* desc:
*
 */

package cluster

import (
	"errors"
	"strings"
	"topiik/internal/config"
	"topiik/internal/util"
	"topiik/node"
)

const clusterInitFailed = "cluster init failed"

/*
* Execute command INIT-CLUSTER
*
 */
func InitCluster(ptnCount int, serverConfig *config.ServerConfig) (ptnIds []string, err error) {
	l.Info().Msg("cluster::ClusterInit start")

	// 0. init cluster
	err = doInit(serverConfig)
	if err != nil {
		return ptnIds, err
	}
	node.InitCluster(controllerInfo.ClusterId)
	nodeStatus.Role = RAFT_LEADER

	/* create partition */
	ptnIds, err = NewPartition(ptnCount)

	if err != nil {
		l.Err(err).Msgf("executor::clusterInit %s", err.Error())
		/* TODO: clean cluster info and partition */
		return ptnIds, err
	}

	// after init, the node default is LEADER, and start to AppendEntries()
	go AppendEntries()
	//ptnUpdCh <- struct{}{} // sync partition to followers
	l.Info().Msg("cluster::ClusterInit end")
	return ptnIds, nil
}

func doInit(serverConfig *config.ServerConfig) error {
	if len(node.GetNodeInfo().ClusterId) > 0 {
		return errors.New("current node already in cluster: " + node.GetNodeInfo().Id)
	}
	// set clusterInfo
	nodeId := node.GetNodeInfo().Id
	controllerInfo.ClusterId = strings.ToLower(util.RandStringRunes(10))

	hostPort, err := util.SplitAddress(serverConfig.Listen)
	if err != nil {
		l.Panic().Msg(err.Error())
	}
	controllerInfo.Nodes[nodeId] = node.NodeSlim{
		Id:    nodeId,
		Addr:  serverConfig.Listen,
		Addr2: hostPort[0] + ":" + hostPort[2],
	}

	// save cluster metadata
	saveControllerInfo()
	return nil
}
