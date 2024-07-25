/***
* author: duan hongxing
* date: 21 Jun 2024
* desc:
 */

package executor

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"topiik/cluster"
	"topiik/internal/command"
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/logger"
	"topiik/resp"
)

/***Command RESponse***/
const (
	RES_OK                   = "OK"
	RES_WRONG_ARG            = "WRONG_ARG"
	RES_WRONG_NUMBER_OF_ARGS = "WRONG_NUM_OF_ARGS"
	RES_DATA_TYPE_NOT_MATCH  = "DATA_TYPE_NOT_MATCH"
	RES_SYNTAX_ERROR         = "SYNTAX_ERR"
	RES_KEY_NOT_EXIST        = "KEY_NOT_EXIST"
	RES_KEY_EXIST_ALREADY    = "KEY_EXIST_ALREADY"

	/*** VOTE response ***/
	RES_ACCEPTED = "A"
	RES_REJECTED = "R"
)

var log = logger.Get()

func Execute(msg []byte, srcAddr string, serverConfig *config.ServerConfig) []byte {
	msgData := msg[4:] // strip the lenght header
	// split msg into [CMD, params]
	//strs := strings.SplitN(strings.TrimLeft(string(msgData), consts.SPACE), consts.SPACE, 2)
	//CMD := strings.ToUpper(strings.TrimSpace(strs[0]))

	icmd, _, err := proto.DecodeHeader(msgData)
	if err != nil {
		log.Err(err)
	}

	if len(msgData) < 2 {
		return resp.ErrorResponse(errors.New(resp.RES_SYNTAX_ERROR))
	}
	var req datatype.Req
	err = json.Unmarshal(msgData[2:], &req) // 2= 1 icmd and 1 ver
	if err != nil {
		log.Err(err).Msg(err.Error())
		return resp.ErrorResponse(err)
	}
	//log.Info().Msgf("aaa %s", req.CMD)
	//CMD := strings.ToUpper(req.CMD)

	//pieces, err := util.SplitCommandLine(string(msgData[2:]))
	//if err != nil {
	//	return resp.ErrorResponse(err)
	//}
	//log.Info().Msg(string(msgData[2:]))
	//log.Info().Msg(strings.Join(pieces, ","))

	pieces := []string{}
	if icmd == command.INIT_CLUSTER_I {
		err := clusterInit(pieces, serverConfig)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(RES_OK, icmd, nil)
	} else if icmd == command.ADD_NODE_I {
		result, err := addNode(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result, icmd, nil)
	} else if icmd == command.SCALE_I {
		result, err := scale(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result, icmd, msg)
	} else if icmd == command.GET_LEADER_ADDR_I {
		log.Info().Msg("get controller address")
		var address string
		if cluster.GetNodeStatus().Role == cluster.RAFT_LEADER { // if is leader, then just return leader's address
			address = serverConfig.Listen
		} else {
			address = cluster.GetNodeStatus().LeaderControllerAddr
		}
		// if not current not controller leader, nor in any cluster, i.e. LeaderControllerAddr is empty
		// then use listen address
		if address == "" {
			address = serverConfig.Listen
		}
		return resp.StringResponse(address, icmd, msg)
	}

	// cluster command
	/*if icmd == command.INIT_CLUSTER {
		pieces := splitParams(strs)
		if len(pieces) < 1 {
			return resp.ErrorResponse(errors.New(RES_SYNTAX_ERROR))
		}
		fmt.Println(pieces)
		if strings.ToUpper(pieces[0]) == "INIT" {
			err := clusterInit(pieces, serverConfig)
			if err != nil {
				return resp.ErrorResponse(err)
			}
			return resp.StringResponse(RES_OK, CMD, msg)
		} else if strings.ToUpper(pieces[0]) == "JOIN" { // CLUSTER JOIN host:port CONTROLLER|WORKER
			if len(pieces) < 3 {
				return resp.ErrorResponse(errors.New(RES_SYNTAX_ERROR))
			}
			result, err := clusterJoin(serverConfig.Listen, pieces[1], pieces[2])
			if err != nil {
				return resp.ErrorResponse(err)
			}
			return resp.StringResponse(result, CMD, msg)
		} else if strings.ToUpper(pieces[0]) == "INFO" {
			clusterInfo()
		}
		return resp.ErrorResponse(errors.New(RES_SYNTAX_ERROR))
	}*/

	// if is Controller, forward to worker(s)
	if cluster.IsNodeController() {
		return cluster.Forward(msg)
	}

	// node must be in a cluster
	if len(cluster.GetNodeInfo().ClusterId) == 0 {
		return resp.ErrorResponse(errors.New("current node not member of cluster"))
	}
	// allow cmd only from Controller Leader, and TODO: allow from Partition Leader
	err = srcFilter(srcAddr)
	if err != nil {
		return resp.ErrorResponse(err)
	}

	if icmd == command.GET_I { // STRING COMMANDS
		result, err := get(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result, icmd, msg)
	} else if icmd == command.SET_I {
		result, err := set(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result, icmd, msg)
	} else if icmd == command.GETM_I {
		result, err := getM(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringArrayResponse(result, icmd, msg)
	} else if icmd == command.SETM_I {
		result, err := setM(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.IntegerResponse(int64(result), icmd, msg)
	} else if icmd == command.INCR_I {
		result, err := incr(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.IntegerResponse(result, icmd, msg)
	} else if icmd == command.LPUSH_I || icmd == command.LPUSHR_I { // LIST COMMANDS
		/***List LPUSH***/
		result, err := pushList(pieces, icmd)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.IntegerResponse(int64(result), icmd, msg)
	} else if icmd == command.LPOP_I || icmd == command.LPOPR_I {
		result, err := popList(pieces, icmd)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringArrayResponse(result, icmd, msg)
	} else if icmd == command.LLEN_I {
		result, err := llen(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.IntegerResponse(int64(result), icmd, msg)
	} else if icmd == command.TTL_I { // KEY COMMANDS
		result, err := ttl(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.IntegerResponse(result, icmd, msg)
	} else if icmd == command.KEYS_I {
		result, err := keys(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringArrayResponse(result, icmd, msg)
	} else {
		fmt.Printf("Invalid cmd: %s\n", icmd)
		return resp.ErrorResponse(errors.New(consts.RES_INVALID_CMD))
	}
}

/***
** Parse the command line after stripped CMD, for commands that KEY is required
** Return:
**	- pieces: [0] the KEY, [1] the request params if have
**	- error: if len(pieces) < 1, means no KEY
**/
func needKEY(cmdKeyParams []string) (pieces []string, err error) {
	if len(cmdKeyParams) < 2 { // first piece is CMD, second is KEY
		return nil, errors.New(RES_SYNTAX_ERROR)
	}
	return strings.SplitN(strings.TrimLeft(cmdKeyParams[1], consts.SPACE), consts.SPACE, 2), nil
}

/***
** Persist command
**
**/
/*
func enqueuePersistentMsg(msg []byte) {
	if memo.MemMap[consts.PERSISTENT_BUF_QUEUE] == nil {
		memo.MemMap[consts.PERSISTENT_BUF_QUEUE] = &datatype.TValue{
			Typ: datatype.V_TYPE_LIST,
			Lst: list.New(),
			Exp: consts.UINT32_MAX,
		}
	}
	memo.MemMap[consts.PERSISTENT_BUF_QUEUE].Lst.PushFront(msg)
}*/

/*** Response json ***/
/*func errorResponse(err error) []byte {
	return response[string](false, err.Error())
}

func successResponse[T any](result T, CMD string, msg []byte) []byte {
	if slices.Contains(needPersistCMD, CMD) {
		enqueuePersistentMsg(msg)
	}
	return response[T](true, result)
}

func response[T any](success bool, response T) []byte {
	b, _ := json.Marshal(&datatype.Response[T]{R: success, M: response})
	return b
}*/

func srcFilter(srcAddr string) error {
	// if node member of cluster
	if len(cluster.GetNodeInfo().ClusterId) > 0 {
		if !cluster.IsNodeController() {
			//fmt.Printf("remote address: %s\n", srcAddr)
			addrSplit, err := util.SplitAddress(srcAddr)
			if err != nil {
				return errors.New(consts.RES_INVLID_OP_ON_WORKER)
			}
			// TOTO: if source host is not Leader's host, also reject
			// if source port is not forward port, also reject
			if addrSplit[1] != cluster.CONTROLLER_FORWORD_PORT {
				fmt.Println(addrSplit[1])
				return errors.New(consts.RES_INVLID_OP_ON_WORKER)
			}
		}
	}
	return nil
	/*if cluster.GetNodeStatus().Role == cluster.RAFT_FOLLOWER {

	}*/
}
