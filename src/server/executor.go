/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"topiik/cluster"
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

	var icmd uint8 // two bytes of command
	if len(msg) >= 1 {
		byteBuf := bytes.NewBuffer([]byte{msg[0]})
		err := binary.Read(byteBuf, binary.LittleEndian, &icmd)
		if err != nil {
			fmt.Println("(err):")
		}
	}

	dataBytes := msg[1:]

	if icmd == cluster.RPC_SYNC_BINLOG {
		/*
		* RPC from worker slave to worker leader
		* To fetching(sync) binary log
		 */

	} else if icmd == cluster.RPC_ADD_NODE {
		/*
		* Client connect to controller leader, and issue ADD-NODE command
		* RPC from controller leader, to add current node to cluster
		 */
		pieces := strings.Split(string(dataBytes), consts.SPACE)
		result, err := addNode(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result)
	} else if icmd == cluster.RPC_VOTE {
		/*
		* RPC from controller leader by request vote
		 */
		cTerm, err := strconv.Atoi(string(dataBytes))
		if err != nil {
			return resp.ErrorResponse(errors.New(RES_SYNTAX_ERROR))
		} else {
			result := vote(cTerm)
			return resp.StringResponse(result)
		}
	} else if icmd == cluster.RPC_APPENDENTRY {
		/*
		* RPC from controller leader by append entry periodic task
		 */
		err := appendEntry(dataBytes, serverConfig)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse("")
	} else if icmd == cluster.RPC_GET_PL {
		/*
		* RPC from workers, to get partition leader addr2
		* for sync data from partition leader
		 */
		pieces := strings.Split(string(dataBytes), consts.SPACE)
		res, err := getPartitionLeader(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(res)
	}

	/* Obsoleted CLUSTER JOIN
	if icmd == CLUSTER_JOIN_ACK {
		pieces := strings.Split(string(dataBytes), consts.SPACE)
		if len(pieces) < 1 {
			return resp.ErrorResponse(errors.New(RES_SYNTAX_ERROR))
		}
		result, err := clusterJoin(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result)
	}*/
	return resp.ErrorResponse(errors.New(consts.RES_INVALID_CMD))
}
