/***
** author: duan hongxing
** data: 13 Jul 2024
** desc:
**
**/

package executor

import (
	"errors"
	"fmt"
	"strconv"
	"topiik/cluster"
	"topiik/internal/config"
)

/***
** Init a Topiik cluster with number of partitions and replicas
** Parameters:
**	- pieces:
**	- serverConfig:
** Syntax: CLUSTER INIT partitions replicas
**/
func clusterInit(pieces []string, serverConfig *config.ServerConfig) error {
	// validate

	// if node already in a cluster, return error
	if len(cluster.GetNodeInfo().ClusterId) > 0 {
		return errors.New("current node already in cluster:" + cluster.GetNodeInfo().ClusterId)
	}

	if len(pieces) != 3 {
		return errors.New(RES_SYNTAX_ERROR)
	}
	partitions, err := strconv.Atoi(pieces[1])
	if err != nil || partitions < 1 {
		fmt.Printf("%s invalid partition number: %s", RES_SYNTAX_ERROR, pieces[1])
		return errors.New(RES_SYNTAX_ERROR)
	}
	replicas, err := strconv.Atoi(pieces[2])
	if err != nil || replicas < 1 {
		fmt.Printf("%s invalid replica number: %s", RES_SYNTAX_ERROR, pieces[2])
		return errors.New(RES_SYNTAX_ERROR)
	}

	return cluster.ClusterInit(uint16(partitions), uint16(replicas), serverConfig)
}
