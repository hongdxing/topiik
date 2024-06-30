/***
** author: duan hongxing
** date: 30 Jun 2024
** desc:
**	Init cluster command implement
**
**/

package executor

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"topiik/cluster"
	"topiik/internal/config"
)

/***
** Init a cluster
** Parameters:
**	- pieces: command line that CMD(CLUSTER INIT) stripped
** Algorithm:
**	1) If any of the node already in cluster, then return error
**	2) If any of the node have data already, then return error
**
** Syntax: CLUSTER INIT --nodes node2:port2 node3:port3
**/
func clusterInit(pieces []string, serverConfig *config.ServerConfig) (err error) {
	if len(pieces) < 3 { //at least --nodes node2:port2 node3:port3
		return errors.New(RES_SYNTAX_ERROR)
	}
	var nodes = []string{}
	if strings.ToUpper(pieces[0]) == "--NODES" {
		nodes = pieces[1:3]
		reg, _ := regexp.Compile(`^.+:\d{4,5}`)
		for _, addr := range nodes {
			if !reg.MatchString(addr) {
				return errors.New(RES_INVALID_ADDR)
			}
		}
		fmt.Println(nodes)
	} else {
		return errors.New(RES_SYNTAX_ERROR)
	}


	if len(nodes) == 2 { // not include current node
		err = cluster.ClusterInit(nodes, serverConfig)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return
	} else if len(nodes) == 5 { // not include current node
		//
		fmt.Println("TODO: 6 nodes")
	}
	return errors.New(RES_SYNTAX_ERROR)
}
