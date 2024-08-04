/***
** author: duan hongxing
** date: 21 Jul 2024
** desc:
**
**/

package executor

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"regexp"
	"strings"
	"topiik/cluster"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

/*
** Run from controller leader, to add new node to cluster
** The target node must already stated, and not joined any cluster yet
** Syntax: ADD-NODE host:port CONTROLLER|WORKER
**
 */
func addNode(req datatype.Req) (result string, err error) {
	pieces, err := util.SplitCommandLine(req.ARGS)
	if len(pieces) != 2 { // must have target address
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	nodeAddr := pieces[0]
	role := strings.ToUpper(pieces[1])
	// validate host
	reg, _ := regexp.Compile(consts.HOST_PATTERN)
	if !reg.MatchString(nodeAddr) {
		return "", errors.New("invalide address")
	}

	// validate CONTROLLER|WORKER
	if strings.ToUpper(role) != cluster.ROLE_CONTROLLER && strings.ToUpper(role) != cluster.ROLE_WORKER {
		return "", errors.New("invalide role, must be either CONTROLLER or WORKER")
	}

	// get controller address2
	addrSplit, _ := util.SplitAddress(nodeAddr)
	nodeAddr2 := addrSplit[0] + ":" + addrSplit[2]
	conn, err := util.PreapareSocketClient(nodeAddr2)
	if err != nil {
		return "", errors.New(resp.RES_NODE_NA)
	}
	defer conn.Close()

	var cmdBytes []byte
	var byteBuf = new(bytes.Buffer) // int to byte buf
	_ = binary.Write(byteBuf, binary.LittleEndian, cluster.RPC_ADD_NODE)
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
	var flag int8
	err = binary.Read(flagBuff, binary.LittleEndian, &flag)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return "", errors.New("add node failed")
	}
	resp := string(buf[resp.RESPONSE_HEADER_SIZE:]) // the node id

	if flag == 1 {
		l.Info().Msgf("Add node succeed:%s", resp)
		cluster.AddNode(resp, nodeAddr, nodeAddr2, role)

	} else {
		l.Err(nil).Msg("Add node failed")
		return "", errors.New(resp)
	}

	return RES_OK, nil
}
