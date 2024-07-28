/***
* author: duan hongxing
* date: 21 Jun 2024
* desc:
 */

package executor

import (
	"encoding/json"
	"errors"
	"slices"
	"topiik/cluster"
	"topiik/internal/command"
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/proto"
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
)

var log = logger.Get()

var PersistenceCh = make(chan []byte)
var persistCmds = []uint8{
	// String
	command.SET_I,
	command.SETM_I,
	command.INCR_I,
	// List
	command.LPUSH_I,
	command.LPUSHR_I,
	command.LPOP_I,
	command.LPOPR_I,

	//command.LPUSHB_I,
	//command.LPUSHRB_I,
	command.DEL_I,
	command.TTL_I, //??
}

func Execute(msg []byte, srcAddr string, serverConfig *config.ServerConfig) (finalRes []byte) {
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

	if icmd == command.INIT_CLUSTER_I {
		err := clusterInit(req, serverConfig)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(RES_OK, icmd)
	} else if icmd == command.ADD_NODE_I {
		result, err := addNode(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result, icmd)
	} else if icmd == command.SCALE_I {
		result, err := scale(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result, icmd)
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
			if len(cluster.GetNodeInfo().ClusterId) == 0 {
				address = serverConfig.Listen
			} else {
				return resp.ErrorResponse(errors.New(resp.RES_NO_LEADER))
			}

		}
		return resp.StringResponse(address, icmd)
	}

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

	finalRes = Execute1(icmd, req)

	if slices.Contains(persistCmds, icmd) {
		PersistenceCh <- msg
	}
	return finalRes
}

/*
* Execute Memory commands
*
 */
func Execute1(icmd uint8, req datatype.Req) (finalRes []byte) {
	pieces := []string{}
	if icmd == command.GET_I { // STRING COMMANDS
		result, err := get(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.StringResponse(result, icmd)
	} else if icmd == command.SET_I {
		result, err := set(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.StringResponse(result, icmd)
	} else if icmd == command.GETM_I {
		result, err := getM(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.StringArrayResponse(result, icmd)
	} else if icmd == command.SETM_I {
		result, err := setM(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.IntegerResponse(int64(result), icmd)
	} else if icmd == command.INCR_I {
		result, err := incr(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.IntegerResponse(result, icmd)
	} else if icmd == command.LPUSH_I || icmd == command.LPUSHR_I { // LIST COMMANDS
		/***List LPUSH***/
		result, err := pushList(pieces, icmd)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.IntegerResponse(int64(result), icmd)
	} else if icmd == command.LPOP_I || icmd == command.LPOPR_I {
		result, err := popList(pieces, icmd)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.StringArrayResponse(result, icmd)
	} else if icmd == command.LLEN_I {
		result, err := llen(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.IntegerResponse(int64(result), icmd)
	} else if icmd == command.TTL_I { // KEY COMMANDS
		result, err := ttl(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.IntegerResponse(result, icmd)
	} else if icmd == command.KEYS_I {
		result, err := keys(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.StringArrayResponse(result, icmd)
	} else {
		log.Err(errors.New("Invalid cmd:" + string(icmd)))
		return resp.ErrorResponse(errors.New(consts.RES_INVALID_CMD))
	}
	return finalRes
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

			/*addrSplit, err := util.SplitAddress(srcAddr)
			if err != nil {
				return errors.New(consts.RES_INVLID_OP_ON_WORKER)
			}*/

			// TOTO: if source host is not Leader's host, also reject
			// if source port is not forward port, also reject
			/* having problem using the same port
			if addrSplit[1] != cluster.CONTROLLER_FORWORD_PORT {
				fmt.Println(addrSplit[1])
				return errors.New(consts.RES_INVLID_OP_ON_WORKER)
			}*/
		}
	}
	return nil
	/*if cluster.GetNodeStatus().Role == cluster.RAFT_FOLLOWER {

	}*/
}
