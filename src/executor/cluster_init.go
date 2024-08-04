/***
** author: duan hongxing
** data: 13 Jul 2024
** desc:
**
**/

package executor

import (
	"errors"
	"topiik/cluster"
	"topiik/internal/config"
	"topiik/internal/datatype"
	"topiik/node"
)

/***
** Init a Topiik cluster with number of partitions and replicas
** Parameters:
**	- pieces:
**	- serverConfig:
** Syntax: INIT-CLUSTER
**/
func clusterInit(req datatype.Req, serverConfig *config.ServerConfig) error {
	// if node already in a cluster, return error
	if len(node.GetNodeInfo().ClusterId) > 0 {
		return errors.New("current node already in cluster:" + node.GetNodeInfo().ClusterId)
	}

	err := cluster.InitCluster(serverConfig)
	if err != nil {
		l.Panic().Msg(err.Error())
	}
	return nil
}
