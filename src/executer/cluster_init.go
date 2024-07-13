/***
** author: duan hongxing
** data: 13 Jul 2024
** desc:
**
**/

package executer

import (
	"errors"
	"fmt"
	"strconv"
	"topiik/cluster"
	"topiik/internal/config"
)

func clusterInit(pieces []string, serverConfig *config.ServerConfig) error {
	// validate
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
