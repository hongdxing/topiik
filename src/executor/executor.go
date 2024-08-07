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
	"topiik/node"
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
	msgBytes := msg[4:] // strip the lenght header

	icmd, _, err := proto.DecodeHeader(msgBytes)
	if err != nil {
		l.Err(err)
	}

	if len(msgBytes) < 2 {
		return resp.ErrorResponse(errors.New(resp.RES_SYNTAX_ERROR))
	}
	var req datatype.Req
	err = json.Unmarshal(msgBytes[2:], &req) // 2= 1 icmd and 1 ver
	if err != nil {
		l.Err(err).Msg(err.Error())
		return resp.ErrorResponse(err)
	}

	if icmd == command.INIT_CLUSTER_I {
		err := clusterInit(req, serverConfig)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(RES_OK)
	} else if icmd == command.ADD_NODE_I {
		result, err := addNode(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result)
	} else if icmd == command.SCALE_I {
		result, err := scale(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result)
	} else if icmd == command.GET_LADDR_I {
		l.Info().Msg("get controller address")
		var address string
		if cluster.GetNodeStatus().Role == cluster.RAFT_LEADER { // if is leader, then just return leader's address
			address = serverConfig.Listen
		} else {
			address = cluster.GetNodeStatus().LeaderControllerAddr
		}
		// if not current not controller leader, nor in any cluster, i.e. LeaderControllerAddr is empty
		// then use listen address
		if address == "" {
			if len(node.GetNodeInfo().ClusterId) == 0 {
				address = serverConfig.Listen
			} else {
				return resp.ErrorResponse(errors.New(resp.RES_NO_LEADER))
			}

		}
		return resp.StringResponse(address)
	}

	// if is Controller, forward to worker(s)
	if cluster.IsNodeController() {
		return forward(icmd, req, msg)
	}

	// node must be in a cluster
	if len(node.GetNodeInfo().ClusterId) == 0 {
		return resp.ErrorResponse(errors.New(resp.RES_NO_CLUSTER))
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

func forward(icmd uint8, req datatype.Req, msg []byte) []byte {
	// special process SETM, because SETM has more than one keys
	if icmd == command.SETM_I {
		if len(req.KEYS) != len(req.VALS) {
			return resp.ErrorResponse(errors.New(resp.RES_SYNTAX_ERROR))
		}
		// split setm to multi set
		for i, key := range req.KEYS {
			reqN := datatype.Req{KEYS: []string{key}, VALS: []string{req.VALS[i]}} // req object
			reqBytesN, _ := json.Marshal(reqN)                                     // req bytes
			msgN, _ := proto.EncodeHeader(command.SET_I, 1)                        // msg header
			msgN = append(msgN, reqBytesN...)                                      // combine msg header and req bytes
			msgN, _ = proto.EncodeB(msgN)                                          // encode msg
			cluster.Forward(key, msgN)
		}
		return resp.StringResponse(resp.RES_OK)
	} else if icmd == command.GETM_I {
		var res []string
		// split setm to multi set
		for _, key := range req.KEYS {
			reqN := datatype.Req{KEYS: []string{key}, VALS: []string{}} // req object
			reqBytesN, _ := json.Marshal(reqN)                          // req bytes
			msgN, _ := proto.EncodeHeader(command.GET_I, 1)             // msg header
			msgN = append(msgN, reqBytesN...)                           // combine msg header and req bytes
			msgN, _ = proto.EncodeB(msgN)                               // encode msg
			resN := cluster.Forward(key, msgN)
			flag := resp.ParseResFlag(resN)
			if flag != 1 {
				res = append(res, "")
			}
			res = append(res, string(resN[resp.RESPONSE_HEADER_SIZE:]))
		}
		return resp.StringArrayResponse(res)
	} else if icmd == command.KEYS_I {
		res := forwardKeys(msg)
		return resp.StringArrayResponse(res)
	}
	key := req.KEYS[0]
	return cluster.Forward(key, msg)
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
		finalRes = resp.StringResponse(result)
	} else if icmd == command.SET_I {
		result, err := set(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.StringResponse(result)
	} else if icmd == command.GETM_I {
		result, err := getM(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.StringArrayResponse(result)
	} else if icmd == command.SETM_I {
		result, err := setM(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.IntegerResponse(int64(result))
	} else if icmd == command.INCR_I {
		result, err := incr(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.IntegerResponse(result)
	} else if icmd == command.LPUSH_I || icmd == command.LPUSHR_I { // LIST COMMANDS
		/***List LPUSH***/
		result, err := pushList(pieces, icmd)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.IntegerResponse(int64(result))
	} else if icmd == command.LPOP_I || icmd == command.LPOPR_I {
		result, err := popList(pieces, icmd)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.StringArrayResponse(result)
	} else if icmd == command.LLEN_I {
		result, err := llen(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.IntegerResponse(int64(result))
	} else if icmd == command.TTL_I { // KEY COMMANDS
		result, err := ttl(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.IntegerResponse(result)
	} else if icmd == command.KEYS_I {
		result, err := keys(req)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		finalRes = resp.StringArrayResponse(result)
	} else {
		l.Err(errors.New("Invalid cmd:" + string(icmd)))
		return resp.ErrorResponse(errors.New(consts.RES_INVALID_CMD))
	}
	return finalRes
}

func srcFilter(srcAddr string) error {
	// if node member of cluster
	if len(node.GetNodeInfo().ClusterId) > 0 {
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
