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
// Param:
//   - addr: current node addr
//
// Syntax: INIT-CLUSTER partition count [controller host:port[,host:port]] worker host:port[,host:port]
func ClusterInit(req datatype.Req, addr string) (err error) {
	pieces := strings.Split(req.Args, consts.SPACE)
	// if node already in a cluster, return error
	if len(node.GetNodeInfo().ClusterId) > 0 {
		return errors.New("current node already in cluster: " + node.GetNodeInfo().ClusterId)
	}

	// at least have paritition and worker
	if len(pieces) != 4 && len(pieces) != 6 {
		return errors.New(resp.RES_SYNTAX_ERROR)
	}

	var argv = make([]string, 2)
	// partition
	argv = pieces[:2]
	if strings.ToLower(string(argv[0])) != "partition" {
		return errors.New(resp.RES_SYNTAX_ERROR)
	}
	var ptnCount int
	ptnCount, err = strconv.Atoi(string(argv[1]))
	if err != nil {
		return err
	}
	if ptnCount <= 0 {
		fmt.Printf("partition: %v", ptnCount)
		return errors.New(resp.RES_SYNTAX_ERROR)
	}

	var (
		controllers []string
		workers     []string
	)
	// controller and worker
	argv = pieces[2:4]
	if argv[0] == "controller" {
		controllers = strings.Split(argv[1], ",")
	} else if argv[0] == "worker" {
		workers = strings.Split(argv[1], ",")
	}

	if len(pieces) == 6 {
		argv = pieces[4:]
		if argv[0] == "controller" {
			controllers = strings.Split(argv[1], ",")
		} else if argv[0] == "worker" {
			workers = strings.Split(argv[1], ",")
		}
	}

	// validate workers
	if len(workers) == 0 || len(workers) < ptnCount {
		return errors.New(resp.RES_NEED_MORE_WORKER)
	}

	controllers = append(controllers, addr)

	// connective check for controllers
	ctlNodeIdAddr, ctlNodeIdAddr2 := checkConnection(controllers)
	if len(ctlNodeIdAddr) != len(controllers) {
		unaccessibleAddr := controllers[len(ctlNodeIdAddr)]
		l.Err(nil).Msgf("Invalid address: %s not accessible", unaccessibleAddr)
		return errors.New(resp.RES_NODE_NA)
	}

	// connective check for workers
	wrkNodeIdAddr, wrkNodeIdAddr2 := checkConnection(workers)
	if len(wrkNodeIdAddr) != len(workers) {
		unaccessibleAddr := workers[len(wrkNodeIdAddr)]
		l.Err(nil).Msgf("Invalid address: %s not accessible", unaccessibleAddr)
		return errors.New(resp.RES_NODE_NA)
	}

	// init cluster
	err = cluster.InitCluster(ctlNodeIdAddr, wrkNodeIdAddr, ptnCount)
	if err != nil {
		l.Err(err).Msgf("executor::clusterInit %s", err.Error())
		/* TODO: clean cluster info if failed */
		return err
	}

	// update controllers
	for ndId, addr2 := range ctlNodeIdAddr2 {
		if ndId == node.GetNodeInfo().Id {
			continue
		}
		rpcAddNode(addr2, node.ROLE_CONTROLLER)
	}

	// update worker
	for _, addr2 := range wrkNodeIdAddr2 {
		rpcAddNode(addr2, node.ROLE_WORKER)
	}

	return nil
}

// Check connectivity of nodes
// Return id->addr and id->addr2 maps
func checkConnection(addrs []string) (map[string]string, map[string]string) {
	var addrMap = make(map[string]string)
	var addr2Map = make(map[string]string)
	for _, addr := range addrs {
		host, _, port2, err := util.SplitAddress2(addr)
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
