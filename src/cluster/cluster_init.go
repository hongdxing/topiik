// author: duan hongxing
// date: 6 July, 2024

package cluster

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
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

	// 1. reshard to assign Slots
	err = ReShard(true)
	if err != nil {
		l.Err(err).Msgf("executor::clusterInit %s", err.Error())
		/* TODO: clean cluster info and partition */
		return err
	}

	// 2. rpc update workers
	for _, group := range workerGroupInfo.Groups {
		for _, nd := range group.Nodes {
			if nd.Id == node.GetNodeInfo().Id {
				continue
			}
			rpcAddNode(nd.Addr2, workerGroupInfo.ClusterId, group.Id)
		}
	}

	// 3. send notification to sync worker group to other workers
	notifyWorkerGroupChanged()

	// 4. after init, start RequestVote
	go RequestVote()

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
	var currentNodeWgId string
	for i := 0; i < ptnCount; i++ {
		addrIdx = 0
		workerGroup := WorkerGroup{Nodes: make(map[string]node.NodeSlim)}
		wgId := strings.ToLower(util.RandStringRunes(10))
		workerGroup.Id = wgId
		workerGroupInfo.Groups[wgId] = &workerGroup
		for ndId, addr := range workers {
			if ndId == node.GetNodeInfo().Id {
				currentNodeWgId = wgId
			}
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
	node.InitCluster(clusterId, currentNodeWgId)

	// save controllerInfo and workerInfo
	saveWorkerGroups()
	return nil
}

func rpcAddNode(addr2 string, clusterId string, grpId string) (string, error) {
	conn, err := util.PreapareSocketClient(addr2)
	if err != nil {
		return "", errors.New(resp.RES_NODE_NA)
	}
	defer conn.Close()

	var buf []byte
	var bbuf = new(bytes.Buffer) // int to byte buf
	_ = binary.Write(bbuf, binary.LittleEndian, consts.RPC_ADD_NODE)
	buf = append(buf, bbuf.Bytes()...)
	//line = clusterid role
	line := clusterId + consts.SPACE + grpId
	buf = append(buf, []byte(line)...)

	// encode
	data, err := proto.EncodeB(buf)
	if err != nil {
		l.Err(err).Msg(err.Error())
	}

	// write
	_, err = conn.Write(data)
	if err != nil {
		l.Err(err).Msg(err.Error())
	}

	// read
	reader := bufio.NewReader(conn)
	buf, err = proto.Decode(reader)
	if err != nil {
		return "", errors.New(resp.RES_NODE_NA)
	}

	flag := resp.ParseResFlag(buf)
	if flag == resp.Success {
		ndId := string(buf[resp.RESPONSE_HEADER_SIZE:]) // the node id
		return ndId, nil
	}
	return "", errors.New(resp.RES_NODE_NA)
}
