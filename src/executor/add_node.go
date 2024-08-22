/*
* author: duan hongxing
* date: 21 Jul 2024
* desc:
*
 */

package executor

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"regexp"
	"topiik/cluster"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

/*
* Add worker to cluster
* Syntax: ADD-WORKER host:port partition {ptnId}
 */
func addWorker(req datatype.Req) (rslt string, err error) {
	pieces, err := util.SplitCommandLine(req.ARGS)
	if len(pieces) != 3 { // must have target address
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	nodeAddr := pieces[0]
	// validate host
	reg, _ := regexp.Compile(consts.HOST_PATTERN)
	if !reg.MatchString(nodeAddr) {
		return "", errors.New("invalide address")
	}

	if pieces[1] != "partition" {
		return rslt, errors.New(RES_SYNTAX_ERROR)
	}

	ptnId := pieces[3]
	if _, ok := cluster.GetPartitionInfo().PtnMap[ptnId]; !ok {
		return rslt, errors.New(resp.RES_INVALID_PARTITION_ID)
	}

	rslt, err = addNode(nodeAddr, cluster.ROLE_WORKER)

	return rslt, err
}

/*
* Add controller to cluster
* Syntax: ADD-CONTROLLER host:port
 */
func addController(req datatype.Req) (rslt string, err error) {
	pieces, err := util.SplitCommandLine(req.ARGS)
	if len(pieces) != 1 { // must have target address
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	nodeAddr := pieces[0]
	// validate host
	reg, _ := regexp.Compile(consts.HOST_PATTERN)
	if !reg.MatchString(nodeAddr) {
		return "", errors.New("invalide address")
	}
	rslt, err = addNode(nodeAddr, cluster.ROLE_CONTROLLER)
	return rslt, err
}

/*
* Run from controller leader, to add new node to cluster
* The target node must already stated, and not joined any cluster yet
* Syntax: ADD-NODE host:port CONTROLLER|WORKER partition xxx
*
 */
func addNode(nodeAddr string, role string) (result string, err error) {
	/* get controller address2 */
	hostPort, _ := util.SplitAddress(nodeAddr)
	nodeAddr2 := hostPort[0] + ":" + hostPort[2]
	conn, err := util.PreapareSocketClient(nodeAddr2)
	if err != nil {
		return "", errors.New(resp.RES_NODE_NA)
	}
	defer conn.Close()

	var cmdBytes []byte
	var byteBuf = new(bytes.Buffer) // int to byte buf
	_ = binary.Write(byteBuf, binary.LittleEndian, consts.RPC_ADD_NODE)
	cmdBytes = append(cmdBytes, byteBuf.Bytes()...)
	//line = clusterid role
	line := node.GetNodeInfo().ClusterId + consts.SPACE + role
	cmdBytes = append(cmdBytes, []byte(line)...)

	data, err := proto.EncodeB(cmdBytes)
	if err != nil {
		l.Err(err).Msg(err.Error())
	}

	// Send
	_, err = conn.Write(data)
	if err != nil {
		l.Err(err).Msg(err.Error())
	}

	reader := bufio.NewReader(conn)
	buf, err := proto.Decode(reader)
	flagByte := buf[4:5]
	flagBuff := bytes.NewBuffer(flagByte)
	var flag resp.RespFlag
	err = binary.Read(flagBuff, binary.LittleEndian, &flag)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return "", errors.New("add node failed")
	}
	res := string(buf[resp.RESPONSE_HEADER_SIZE:]) // the node id

	if flag == resp.Success {
		l.Info().Msgf("executor::addNode succeed:%s", res)
		cluster.AddNode(res, nodeAddr, nodeAddr2, role)

	} else {
		l.Err(nil).Msg("executor::addNode failed")
		return "", errors.New(res)
	}

	return RES_OK, nil
}
