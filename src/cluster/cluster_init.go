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
func InitCluster(workers map[string]string, ptnCount int) (err error) {
	l.Info().Msg("cluster::ClusterInit start")

	// 0. init cluster
	// generate cluster id, set workers
	err = doInit(workers, ptnCount)
	if err != nil {
		return err
	}

	// 3. reshard to assign Slots
	err = ReShard(true)
	if err != nil {
		l.Err(err).Msgf("executor::clusterInit %s", err.Error())
		/* TODO: clean cluster info and partition */
		return err
	}

	// 4. send notification to sync meta data to other controller(s) and worker(s)
	notifyWorkerGroupChanged()

	// 5. after init, the node default is LEADER, and start to AppendEntries()
	//go AppendEntries()
	go RequestVote()
	//ptnUpdCh <- struct{}{} // sync partition to followers
	l.Info().Msg("cluster::ClusterInit end")
	return nil
}

func doInit(workers map[string]string, ptnCount int) error {
	if len(node.GetNodeInfo().ClusterId) > 0 {
		return errors.New("current node already in cluster: " + node.GetNodeInfo().Id)
	}

	// generate cluster id
	clusterId := strings.ToLower(util.RandStringRunes(10))
	workerGroupInfo.ClusterId = clusterId

	// set controllerInfo
	var addrIdx = 0
	for i := 0; i < ptnCount; i++ {
		addrIdx = 0
		workerGroup := WorkerGroup{Nodes: make(map[string]node.NodeSlim)}
		wgId := strings.ToLower(util.RandStringRunes(10))
		workerGroup.Id = wgId
		workerGroupInfo.Groups[wgId] = &workerGroup
		for ndId, addr := range workers {
			if addrIdx%int(ptnCount) == i {
				host, _, port2, _ := util.SplitAddress2(addr)
				workerGroup.Nodes[ndId] = node.NodeSlim{Id: ndId, Addr: addr, Addr2: host + ":" + port2}
				// set first to leader
				if len(workerGroup.Nodes) == 1 {
					workerGroup.LeaderNodeId = ndId
				}
			}
			addrIdx++
		}
	}

	// update current(controller) node
	node.InitCluster(clusterId)

	// save controllerInfo and workerInfo
	saveWorkerGroups()
	return nil
}
