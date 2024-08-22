/*
* author: duan hongxing
* data: 13 Jul 2024
* desc:
*
 */

package executor

import (
	"errors"
	"strconv"
	"topiik/cluster"
	"topiik/internal/config"
	"topiik/internal/datatype"
	"topiik/node"
	"topiik/resp"
)

/*
* Init a Topiik cluster
* Parameters:
*	- serverConfig
* Syntax: INIT-CLUSTER partition count
 */
func clusterInit(req datatype.Req, serverConfig *config.ServerConfig) (ptnIds []string, err error) {
	/* if node already in a cluster, return error */
	if len(node.GetNodeInfo().ClusterId) > 0 {
		return nil, errors.New("current node already in cluster: " + node.GetNodeInfo().ClusterId)
	}

	/* validate paritition */
	if len(req.ARGS) == 0 {
		return nil, errors.New(resp.RES_SYNTAX_ERROR)
	}
	var ptnCount int
	ptnCount, err = strconv.Atoi(req.ARGS)
	if err != nil {
		return nil, err
	}
	if ptnCount <= 0 {
		return nil, errors.New(resp.RES_SYNTAX_ERROR)
	}

	/* init cluster */
	l.Info().Msg("before")
	ptnIds, err = cluster.InitCluster(ptnCount, serverConfig)
	l.Info().Msg("after")
	if err != nil {
		l.Err(err).Msgf("executor::clusterInit %s", err.Error())
		/* TODO: clean cluster info if failed */
		return ptnIds, err
	}

	return ptnIds, nil
}
