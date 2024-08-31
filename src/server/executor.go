/*
* author: duan hongxing
* date: 3 Jul 2024
* desc:
*
 */

package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"topiik/executor"
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/persistence"
	"topiik/resp"
)

/*
* msg:
*	- msg[:1] 	= command
*	- msg[1:]	= data
*
 */
func Execute(msg []byte, serverConfig *config.ServerConfig) (result []byte) {
	// TODO: remove
	//strs := strings.SplitN(strings.TrimLeft(string(msg[4:]), consts.SPACE), consts.SPACE, 2)
	//CMD := strings.ToUpper(strings.TrimSpace(strs[0]))
	if len(msg) < 4 {
		return resp.ErrResponse(errors.New(resp.RES_NIL))
	}
	msg = msg[4:]

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
		/*if len(dataBytes) < consts.NODE_ID_LEN {
			return resp.ErrResponse(errors.New(resp.RES_SYNTAX_ERROR))
		}
		followerId := string(dataBytes[:consts.NODE_ID_LEN])
		var seq int64
		byteBuf := bytes.NewBuffer(dataBytes[consts.NODE_ID_LEN:])
		err := binary.Read(byteBuf, binary.LittleEndian, &seq)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return resp.ErrResponse(errors.New(resp.RES_SYNTAX_ERROR))
		}
		res := persistence.Fetch(followerId, seq)
		*/

		/*
		* RPC from Partition Leader sync binlog to follower
		 */
		res, err := persistence.ReceiveBinlog(dataBytes, executor.Execute1)
		if err != nil {
			return resp.ErrResponse(err)
		}
		return resp.IntResponse(res)
	} else if icmd == consts.RPC_GET_BLSEQ {
		/*
		* controller get worker node binlog seq, to elect new partition leader
		 */
		seq := persistence.GetBLSeq()
		return resp.IntResponse(seq)
	} else if icmd == consts.RPC_ADD_NODE {
		/*
		* Client connect to controller leader, and issue ADD-NODE command
		* RPC from controller leader, to add current node to cluster
		 */
		pieces := strings.Split(string(dataBytes), consts.SPACE)
		result, err := addNode(pieces)
		if err != nil {
			return resp.ErrResponse(err)
		}
		return resp.StrResponse(result)
	} else if icmd == consts.RPC_VOTE {
		/*
		* RPC from controller leader by request vote
		 */
		cTerm, err := strconv.Atoi(string(dataBytes))
		if err != nil {
			return resp.ErrResponse(errors.New(resp.RES_SYNTAX_ERROR))
		} else {
			result := vote(cTerm)
			return resp.StrResponse(result)
		}
	} else if icmd == consts.RPC_APPENDENTRY {
		/*
		* RPC from controller leader by append entry periodic task
		 */
		err := appendEntry(dataBytes, serverConfig)
		if err != nil {
			return resp.ErrResponse(err)
		}
		return resp.StrResponse("")
	} else if icmd == consts.RPC_GET_PL {
		/*
		* RPC from workers, to get partition leader addr2
		* for sync data from partition leader
		 */
		pieces := strings.Split(string(dataBytes), consts.SPACE)
		res, err := getPartitionLeader(pieces)
		if err != nil {
			return resp.ErrResponse(err)
		}
		return resp.StrResponse(res)
	}

	/* Obsoleted CLUSTER JOIN
	if icmd == CLUSTER_JOIN_ACK {
		pieces := strings.Split(string(dataBytes), consts.SPACE)
		if len(pieces) < 1 {
			return resp.ErrResponse(errors.New(RES_SYNTAX_ERROR))
		}
		result, err := clusterJoin(pieces)
		if err != nil {
			return resp.ErrResponse(err)
		}
		return resp.StrResponse(result)
	}*/
	return resp.ErrResponse(errors.New(consts.RES_INVALID_CMD))
}
