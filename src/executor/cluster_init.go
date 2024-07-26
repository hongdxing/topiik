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
)

/***
** Init a Topiik cluster with number of partitions and replicas
** Parameters:
**	- pieces:
**	- serverConfig:
** Syntax: INIT-CLUSTER
**/
func clusterInit(req datatype.Req, serverConfig *config.ServerConfig) error {
	// validate

	// if node already in a cluster, return error
	if len(cluster.GetNodeInfo().ClusterId) > 0 {
		return errors.New("current node already in cluster:" + cluster.GetNodeInfo().ClusterId)
	}

	/*if len(pieces) != 2 {
		return errors.New(RES_SYNTAX_ERROR)
	}
	partitionIndex := 0
	replicaIndex := 1
	partitions, err := strconv.Atoi(pieces[partitionIndex])
	if err != nil || partitions < 1 {
		log.Err(err).Msgf("%s invalid partition number: %s", RES_SYNTAX_ERROR, pieces[partitionIndex])
		return errors.New(RES_SYNTAX_ERROR)
	}
	replicas, err := strconv.Atoi(pieces[replicaIndex])
	if err != nil || replicas < 1 {
		log.Err(err).Msgf("%s invalid replica number: %s", RES_SYNTAX_ERROR, pieces[replicaIndex])
		return errors.New(RES_SYNTAX_ERROR)
	}*/

	return cluster.ClusterInit(serverConfig)
}
