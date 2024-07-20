/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package cluster

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/resp"
)

const (
	RES_SYNTAX_ERROR = "SYNTAX_ERR"
)

/*
** msg:
**	- msg[:2] 	= command
**	- msg[2:]	= data
**
 */
func Execute(msg []byte, serverConfig *config.ServerConfig) (result []byte) {
	// TODO: remove
	//strs := strings.SplitN(strings.TrimLeft(string(msg[4:]), consts.SPACE), consts.SPACE, 2)
	//CMD := strings.ToUpper(strings.TrimSpace(strs[0]))
	CMD := ""

	var icmd int16 // two bytes of command
	if len(msg) >= 2 {
		cmdBytes := msg[:2]
		byteBuf := bytes.NewBuffer(cmdBytes)
		err := binary.Read(byteBuf, binary.LittleEndian, &icmd)
		if err != nil {
			fmt.Println("(err):")
		}
	}

	dataBytes := msg[2:]

	if icmd == ClusterCmdMap[CLUSTER_JOIN_ACK] {
		pieces := strings.Split(string(dataBytes), consts.SPACE)
		if len(pieces) < 1 {
			return resp.ErrorResponse(errors.New(RES_SYNTAX_ERROR))
		}
		result, err := clusterJoin(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result, CMD, msg)
	} else if icmd == ClusterCmdMap[RPC_VOTE] {
		cTerm, err := strconv.Atoi(string(dataBytes))
		if err != nil {
			return resp.ErrorResponse(errors.New(RES_SYNTAX_ERROR))
		} else {
			result := vote(cTerm)
			return resp.StringResponse(result, CMD, msg)
		}
	} else if icmd == ClusterCmdMap[RPC_APPENDENTRY] {
		err := appendEntry(dataBytes, serverConfig)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse("", CMD, msg)
	}
	return resp.ErrorResponse(errors.New(consts.RES_INVALID_CMD))
}
