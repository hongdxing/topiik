/*
* author: duan hongxing
* data: 13 Jul 2024
* desc:
*
 */

package executor

import (
	"errors"
	"topiik/cluster"
	"topiik/internal/config"
	"topiik/node"
)

/*
* Init a Topiik cluster
* Parameters:
*	- serverConfig
* Syntax: INIT-CLUSTER
 */
func clusterInit(serverConfig *config.ServerConfig) error {
	// if node already in a cluster, return error
	if len(node.GetNodeInfo().ClusterId) > 0 {
		return errors.New("current node already in cluster: " + node.GetNodeInfo().ClusterId)
	}

	err := cluster.InitCluster(serverConfig)
	if err != nil {
		l.Panic().Msg(err.Error())
	}
	return nil
}
