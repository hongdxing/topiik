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
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/persistence"
	"topiik/resp"
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

	if icmd == consts.RPC_SYNC_BINLOG {
		/*
		* RPC from worker slave to worker leader
		* To fetching(sync) binary log
		* The follow send it's binlogSeq to Leader
		 */
		if len(dataBytes) < consts.NODE_ID_LEN {
			return resp.ErrorResponse(errors.New(resp.RES_SYNTAX_ERROR))
		}
		followerId := string(dataBytes[:consts.NODE_ID_LEN])
		var seq int64
		byteBuf := bytes.NewBuffer(dataBytes[consts.NODE_ID_LEN:])
		err := binary.Read(byteBuf, binary.LittleEndian, &seq)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return resp.ErrorResponse(errors.New(resp.RES_SYNTAX_ERROR))
		}
		res := persistence.Fetch(followerId, seq)
		return resp.StringResponse(string(res))
	} else if icmd == consts.RPC_ADD_NODE {
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
	} else if icmd == consts.RPC_VOTE {
		/*
		* RPC from controller leader by request vote
		 */
		cTerm, err := strconv.Atoi(string(dataBytes))
		if err != nil {
			return resp.ErrorResponse(errors.New(resp.RES_SYNTAX_ERROR))
		} else {
			result := vote(cTerm)
			return resp.StringResponse(result)
		}
	} else if icmd == consts.RPC_APPENDENTRY {
		/*
		* RPC from controller leader by append entry periodic task
		 */
		err := appendEntry(dataBytes, serverConfig)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse("")
	} else if icmd == consts.RPC_GET_PL {
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
