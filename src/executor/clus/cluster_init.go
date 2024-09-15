//author: Duan Hongxing
//data: 13 Jul, 2024

package clus

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"topiik/cluster"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

// Init a Topiik cluster
// Params:
//   - addr: current node addr
func ClusterInit(req datatype.Req, persistors []string) (err error) {
	pieces := strings.Split(req.Args, consts.SPACE)
	l.Info().Msg(req.Args)
	if len(pieces) != 2 {
		return errors.New(resp.RES_SYNTAX_ERROR)
	}
	// if node already in a cluster, return error
	if len(node.GetNodeInfo().ClusterId) > 0 {
		return errors.New("current node already in cluster: " + node.GetNodeInfo().ClusterId)
	}

	workers := strings.Split(pieces[0], ",")
	ptnCount, err := strconv.Atoi(pieces[1])
	if err != nil {
		return err
	}

	// validate persistors
	if len(persistors) == 0 || len(persistors) < ptnCount {
		return errors.New(resp.RES_NO_PERSISTOR)
	}

	// connective check for workers
	ctlNodeIdAddr, _ := checkConnection(workers)
	if len(ctlNodeIdAddr) != len(workers) {
		unaccessibleAddr := workers[len(ctlNodeIdAddr)]
		l.Err(nil).Msgf("Invalid address: %s not accessible", unaccessibleAddr)
		return errors.New(resp.RES_NODE_NA)
	}
	fmt.Println(ctlNodeIdAddr)

	// connective check for persistors
	pstNodeIdAddr, _ := checkConnection(persistors)
	if len(pstNodeIdAddr) != len(persistors) {
		unaccessibleAddr := persistors[len(pstNodeIdAddr)]
		l.Err(nil).Msgf("Invalid address: %s not accessible", unaccessibleAddr)
		return errors.New(resp.RES_NODE_NA)
	}

	// init cluster
	err = cluster.InitCluster(ctlNodeIdAddr, ptnCount)
	if err != nil {
		l.Err(err).Msgf("executor::clusterInit %s", err.Error())
		/* TODO: clean cluster info if failed */
		return err
	}

	// update worker
	//for _, addr2 := range pstNodeIdAddr2 {
	//	rpcAddNode(addr2, node.ROLE_WORKER)
	//}

	return nil
}

// Check connectivity of nodes
// Return id->addr and id->addr2 maps
func checkConnection(addrs []string) (map[string]string, map[string]string) {
	var addrMap = make(map[string]string)
	var addr2Map = make(map[string]string)
	for _, addr := range addrs {
		host, _, port2, err := util.SplitAddress2(strings.TrimSpace(addr))
		if err != nil {
			l.Err(err).Msg(err.Error())
			break
		}

		addr2 := host + ":" + port2
		conn, err := util.PreapareSocketClient(addr2)
		if err != nil {
			l.Err(err).Msg(err.Error())
			break
		}
		defer conn.Close()

		// Prepare buf
		var buf []byte
		bbuf := new(bytes.Buffer)
		binary.Write(bbuf, binary.LittleEndian, consts.RPC_TEST_CONN)
		buf = append(buf, bbuf.Bytes()...)
		buf, err = proto.EncodeB(buf)
		if err != nil {
			break
		}
		// Write
		_, err = conn.Write(buf)
		if err != nil {
			l.Err(err).Msgf("cluster_init::checkConnection write %s", err.Error())
			break
		}
		// Read
		reader := bufio.NewReader(conn)
		res, err := proto.Decode(reader)
		if err != nil {
			l.Err(err).Msgf("cluster_init::checkConnection decode %s", err.Error())
			break
		}

		// Flag
		flag := resp.ParseResFlag(res)

		if flag == resp.Success {
			if len(res) > resp.RESPONSE_HEADER_SIZE {
				ndId := string(res[resp.RESPONSE_HEADER_SIZE:])
				if err != nil {
					l.Err(err).Msgf("cluster_init::checkConnection read %s", err.Error())
					break
				}
				addrMap[ndId] = addr
				addr2Map[ndId] = addr2
			} else {
				l.Warn().Msgf("cluster_init::checkConnection failed")
			}
		} else {
			break
		}
	}

	return addrMap, addr2Map
}
