// author: duan hongxing
// date: 6 July, 2024

package cluster

import (
	"errors"
	"strings"
	"topiik/internal/util"
	"topiik/node"
)

// Execute command INIT-CLUSTER
func InitCluster(controllers map[string]string, ptnCount int) (err error) {
	l.Info().Msg("cluster::ClusterInit start")

	// 0. init cluster
	// generate cluster id, set controllers, set workers
	err = doInit(controllers)
	if err != nil {
		return err
	}

	//var ptnIdx = 0
	for ptnIdx := range ptnCount {
		var wrkIdx = 0
		var nodes = make(map[string]string)
		for ndId, addr := range controllers {
			if wrkIdx%ptnCount == ptnIdx {
				nodes[ndId] = addr
			}
			wrkIdx++
		}
		addPartition(nodes)
		ptnIdx++
	}

	// 3. reshard to assign Slots
	err = ReShard()
	if err != nil {
		l.Err(err).Msgf("executor::clusterInit %s", err.Error())
		/* TODO: clean cluster info and partition */
		return err
	}

	// 4. send notification to sync meta data to other controller(s) and worker(s)
	notifyControllerChanged()
	notifyPtnChanged()

	// 5. after init, the node default is LEADER, and start to AppendEntries()
	go AppendEntries()
	//ptnUpdCh <- struct{}{} // sync partition to followers
	l.Info().Msg("cluster::ClusterInit end")
	return nil
}

func doInit(controllers map[string]string) error {
	if len(node.GetNodeInfo().ClusterId) > 0 {
		return errors.New("current node already in cluster: " + node.GetNodeInfo().Id)
	}

	// generate cluster id
	clusterId := strings.ToLower(util.RandStringRunes(10))

	// set controllerInfo
	controllerInfo.ClusterId = clusterId
	for ndId, addr := range controllers {
		host, _, port2, err := util.SplitAddress2(addr)
		if err != nil {
			l.Panic().Msg(err.Error())
		}
		controllerInfo.Nodes[ndId] = node.NodeSlim{
			Id:    ndId,
			Addr:  addr,
			Addr2: host + ":" + port2,
		}
	}

	// update current(controller) node
	node.InitCluster(clusterId)

	// set current Role to Raft Leader
	nodeStatus.Role = RAFT_LEADER

	// save controllerInfo and workerInfo
	saveControllerInfo()
	return nil
}
