//author: Duan Hongxing
//date: 21 Jul, 2024

package clus

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"regexp"
	"strings"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

// Add worker to cluster
// Syntax: ADD-WORKER host:port partition {ptnId}
func AddWorker(req datatype.Req) (ndId string, err error) {
	pieces, err := util.SplitCommandLine(req.Args)
	if err != nil {
		return "", errors.New(resp.RES_SYNTAX_ERROR)
	}
	if len(pieces) != 3 { // must have target address
		return "", errors.New(resp.RES_SYNTAX_ERROR)
	}
	nodeAddr := pieces[0]
	// validate address
	reg, _ := regexp.Compile(consts.HOST_PATTERN)
	if !reg.MatchString(nodeAddr) {
		return "", errors.New("invalide address")
	}

	if strings.ToUpper(pieces[1]) != "PARTITION" {
		return ndId, errors.New(resp.RES_SYNTAX_ERROR)
	}

	//ndId, err = addNode(nodeAddr, config.ROLE_WORKER, ptnId)
	//if err != nil {
	//	return ndId, err
	//}
	return ndId, err
}

// Run from controller leader, to add new node to cluster
// The target node must already started, and not joined any cluster yet
// This method may trigger via INIT-CLUSTER, ADD-CONTROLLER or ADD-WORKER
func addNode(nodeAddr string, role string, ptnId string) (result string, err error) {
	// get controller address2
	hostPort, _ := util.SplitAddress(nodeAddr)
	nodeAddr2 := hostPort[0] + ":" + hostPort[2]
	conn, err := util.PreapareSocketClient(nodeAddr2)
	if err != nil {
		return "", errors.New(resp.RES_NODE_NA)
	}
	defer conn.Close()

	var cmdBytes []byte
	var bbuf = new(bytes.Buffer) // int to byte buf
	_ = binary.Write(bbuf, binary.LittleEndian, consts.RPC_ADD_NODE)
	cmdBytes = append(cmdBytes, bbuf.Bytes()...)
	//line = clusterid role
	line := node.GetNodeInfo().ClusterId + consts.SPACE + role
	cmdBytes = append(cmdBytes, []byte(line)...)

	// encode
	data, err := proto.EncodeB(cmdBytes)
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
	buf, err := proto.Decode(reader)

	flag := resp.ParseResFlag(buf)
	ndId := string(buf[resp.RESPONSE_HEADER_SIZE:]) // the node id

	if flag == resp.Success {
		l.Info().Msgf("executor::addNode succeed: %s", ndId)
		//cluster.AddNode(ndId, nodeAddr, nodeAddr2, role, ptnId)

	} else {
		l.Err(nil).Msg("executor::addNode failed")
		return "", errors.New(ndId)
	}

	return ndId, nil
}

func rpcAddNode(addr2 string, role string) (string, error) {
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
	line := node.GetNodeInfo().ClusterId + consts.SPACE + role
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
