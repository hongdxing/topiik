/***
** author: duan hongxing
** data: 4 Jul 2024
** desc:
**
**/

package executor

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"topiik/cluster"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

/***
** Join to CCSS cluster
**
**
** Syntax:
**	CLUSTER JOIN host:port CONTROLLER|WORKER
**/
func clusterJoin(myAddr string, controllerAddr string, role string) (result string, err error) {
	// validate host
	reg, _ := regexp.Compile(consts.HOST_PATTERN)
	if !reg.MatchString(myAddr) || !reg.MatchString(controllerAddr) {
		return "", errors.New("invalide address")
	}

	// validate CONTROLLER|WORKER
	if strings.ToUpper(role) != cluster.ROLE_CONTROLLER && strings.ToUpper(role) != cluster.ROLE_WORKER {
		return "", errors.New("invalide role, must be either CONTROLLER or WORKER")
	}
	nodeId := node.GetNodeInfo().Id

	// get controller address2
	addrSplit, _ := util.SplitAddress(controllerAddr)

	conn, err := util.PreapareSocketClient(addrSplit[0] + ":" + addrSplit[2])
	if err != nil {
		return "", errors.New("join to cluster failed, please check whether captial node still alive and try again")
	}
	defer conn.Close()

	// CLUSTER JOIN_ACK nodeId addr role
	//line := cluster.CLUSTER_JOIN_ACK + consts.SPACE + nodeId + consts.SPACE + myAddr + consts.SPACE + role

	var cmdBytes []byte
	var byteBuf = new(bytes.Buffer) // int to byte buf
	_ = binary.Write(byteBuf, binary.LittleEndian, consts.CLUSTER_JOIN_ACK)
	cmdBytes = append(cmdBytes, byteBuf.Bytes()...)
	line := nodeId + consts.SPACE + myAddr + consts.SPACE + role
	cmdBytes = append(cmdBytes, []byte(line)...)

	data, err := proto.EncodeB(cmdBytes)
	if err != nil {
		fmt.Println(err)
	}

	// Send
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println(err)
	}

	reader := bufio.NewReader(conn)
	buf, err := proto.Decode(reader)
	flagByte := buf[4:5]
	flagBuff := bytes.NewBuffer(flagByte)
	var flag int8
	err = binary.Read(flagBuff, binary.LittleEndian, &flag)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("join to cluster failed")
	}
	resp := string(buf[resp.RESPONSE_HEADER_SIZE:])

	if flag == 1 {

		fmt.Printf("join cluster:%s\n", resp)
		//err = cluster.UpdateNodeClusterId(resp)
		if err != nil {
			fmt.Println(err)
			return "", errors.New("join cluster failed")
		}
		// if join controller succeed, will start to RequestVote
		if strings.ToUpper(role) == cluster.ROLE_CONTROLLER {
			go cluster.RequestVote()
		}
	} else {
		fmt.Println("join cluster failed")
		return "", errors.New(resp)
	}

	return RES_OK, nil
}
