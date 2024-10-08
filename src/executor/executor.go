// author: Duan Hongxing
// date: 21 Jun, 2024

package executor

import (
	"encoding/json"
	"errors"
	"fmt"
	"topiik/cluster"
	"topiik/executor/clus"
	"topiik/executor/keyy"
	"topiik/executor/list"
	"topiik/executor/shared"
	"topiik/executor/str"
	"topiik/internal/command"
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/proto"
	"topiik/node"
	"topiik/resp"
)

// var PersistenceCh = make(chan []byte)

func Execute(msg []byte, srcAddr string, serverConfig *config.ServerConfig) (finalRes []byte) {
	msgBytes := msg[4:] // strip the lenght header

	icmd, _, err := proto.DecodeHeader(msgBytes)
	if err != nil {
		l.Err(err)
	}

	if len(msgBytes) < 2 {
		return resp.ErrResponse(errors.New(resp.RES_SYNTAX_ERROR))
	}
	//l.Info().Msgf("%s", string(msgBytes[2:]))
	var req datatype.Req
	err = json.Unmarshal(msgBytes[2:], &req) // 2= 1 icmd and 1 ver
	if err != nil {
		l.Err(err).Msg(err.Error())
		return resp.ErrResponse(err)
	}

	if icmd == command.CREATE_CLUSTER_I {
		err := clus.ClusterInit(req, serverConfig.Persistors)
		if err != nil {
			return resp.ErrResponse(err)
		}
		return resp.StrResponse(resp.RES_OK)
	} else if icmd == command.SHOW_I {
		rslt := clus.Show(req)
		return resp.StrResponse(rslt)
	} else if icmd == command.ADD_WORKER_I {
		_, err := clus.AddWorker(req)
		if err != nil {
			return resp.ErrResponse(err)
		}
		return resp.StrResponse(resp.RES_OK)
	} else if icmd == command.REMOVE_NODE_I {
		err := clus.RemoveNode(req)
		if err != nil {
			return resp.ErrResponse(err)
		}
		return resp.StrResponse(resp.RES_OK)
	} else if icmd == command.NEW_PARTITION_I {
		ptnId, err := clus.NewPartition(req)
		if err != nil {
			return resp.ErrResponse(err)
		}
		return resp.StrResponse(ptnId)
	} else if icmd == command.RESHARD_I {
		result, err := clus.Reshard(req)
		if err != nil {
			return resp.ErrResponse(err)
		}
		return resp.StrResponse(result)
	} else if icmd == command.GET_CTLADDR_I {
		l.Info().Msg("get leader address")

		// allow connection ONLY when node in a cluster
		if node.GetNodeInfo().ClusterId == "" {
			return resp.ErrResponse(fmt.Errorf("%s %s", resp.RES_NO_CLUSTER, "please create cluster"))
		}

		var address string
		if cluster.GetNodeStatus().RaftRole == cluster.RAFT_LEADER { // if is leader, then just return leader's address
			address = serverConfig.Listen
		} else {
			leader := cluster.GetPtnLeader(node.GetNodeInfo().Id)
			address = leader.Addr
		}
		// if not current not controller leader, nor in any cluster, i.e. LeaderControllerAddr is empty
		// then use listen address
		if address == "" {
			return resp.ErrResponse(fmt.Errorf("%s %s", resp.RES_NO_CLUSTER, "please create cluster"))
			/*
				if len(node.GetNodeInfo().ClusterId) == 0 {
					address = serverConfig.Listen
				} else {
					return resp.ErrResponse(errors.New(resp.RES_NO_CTL))
				}
			*/
		}
		return resp.StrResponse(address)
	}

	// if is Controller, forward to worker(s)
	//if node.IsController() {
	//	return forward(icmd, req, msg)
	//}

	// node must be in a cluster
	if len(node.GetNodeInfo().ClusterId) == 0 {
		return resp.ErrResponse(errors.New(resp.RES_NO_CLUSTER))
	}

	//finalRes = Execute1(icmd, req)
	finalRes = forward(icmd, req, msg)

	return finalRes
}

func forward(icmd uint8, req datatype.Req, msg []byte) []byte {
	// special process commands that could have more than one key
	// 1) SETM
	// 2) GETM
	// 3) DEL
	// 4) EXISTS
	if icmd == command.SETM_I {
		if len(req.Keys) != len(req.Vals) {
			return resp.ErrResponse(errors.New(resp.RES_SYNTAX_ERROR))
		}
		// split setm to multi set
		for i, key := range req.Keys {
			reqN := datatype.Req{Keys: datatype.Abytes{key}, Vals: datatype.Abytes{req.Vals[i]}} // req object
			reqBytesN, _ := json.Marshal(reqN)                                                   // req bytes
			msgN, _ := proto.EncodeHeader(command.SET_I, 1)                                      // msg header
			msgN = append(msgN, reqBytesN...)                                                    // combine msg header and req bytes
			msgN, _ = proto.EncodeB(msgN)                                                        // encode msg
			executeOrForward(key, command.SET_I, reqN, msgN)
		}
		return resp.StrResponse(resp.RES_OK)
	} else if icmd == command.GETM_I {
		var res []string
		// split setm to multi set
		for _, key := range req.Keys {
			reqN := datatype.Req{Keys: datatype.Abytes{key}, Vals: datatype.Abytes{}} // req object
			reqBytesN, _ := json.Marshal(reqN)                                        // req bytes
			msgN, _ := proto.EncodeHeader(command.GET_I, 1)                           // msg header
			msgN = append(msgN, reqBytesN...)                                         // combine msg header and req bytes
			msgN, _ = proto.EncodeB(msgN)                                             // encode msg
			resN := executeOrForward(key, command.GET_I, reqN, msgN)
			flag := resp.ParseResFlag(resN)
			if flag != resp.Success {
				res = append(res, "")
			}
			res = append(res, string(resN[resp.RESPONSE_HEADER_SIZE:]))
		}
		return resp.StrArrResponse(res)
	} else if icmd == command.DEL_I {
		rslt := keyy.ForwardDel(Execute1, req, msg)
		return resp.IntResponse(rslt)
	} else if icmd == command.EXISTS_I {
		rslt := keyy.ForwardExists(Execute1, req, msg, len(req.Keys))
		return resp.StrArrResponse(rslt)
	} else if icmd == command.KEYS_I {
		res := keyy.ForwardKeys(Execute1, req, msg)
		return resp.StrArrResponse(res)
	}

	// the key should not empty or space
	if len(req.Keys) <= 0 || len(req.Keys[0]) == 0 {
		return resp.ErrResponse(errors.New(resp.RES_EMPTY_KEY))
	}
	key := req.Keys[0]
	return executeOrForward(key, icmd, req, msg)
}

func executeOrForward(key []byte, icmd uint8, req datatype.Req, msg []byte) []byte {
	targetWorker, err := shared.GetLeaderNode(key)
	if err != nil {
		return resp.ErrResponse(err)
	}
	return shared.ExecuteOrForward(targetWorker, Execute1, icmd, req, msg)
}

// Execute Memory commands
func Execute1(icmd uint8, req datatype.Req) (finalRes []byte, err error) {
	if icmd == command.GET_I { /*** STRING COMMANDS ***/
		result, err := str.Get(req)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.StrResponse(result)
	} else if icmd == command.SET_I {
		result, err := str.Set(req)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.StrResponse(result)
	} else if icmd == command.GETM_I {
		result, err := str.GetM(req)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.StrArrResponse(result)
	} else if icmd == command.SETM_I {
		result, err := str.SetM(req)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.IntResponse(int64(result))
	} else if icmd == command.INCR_I {
		result, err := str.Incr(req)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.IntResponse(result)
	} else if icmd == command.LPUSH_I || icmd == command.LPUSHR_I { // LIST COMMANDS
		result, err := list.Push(req, icmd)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.IntResponse(int64(result))
	} else if icmd == command.LPOP_I || icmd == command.LPOPR_I {
		rslt, err := list.Pop(req, icmd)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		// finalRes = resp.StrArrResponse(result)
		finalRes = resp.StrArrResponse(rslt)
	} else if icmd == command.LLEN_I {
		result, err := list.Len(req)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.IntResponse(int64(result))
	} else if icmd == command.LSLICE_I {
		rslt, err := list.Slice(req)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.StrArrResponse(rslt)
	} else if icmd == command.LSET_I {
		rslt, err := list.Set(req)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.StrResponse(rslt)
	} else if icmd == command.TTL_I { // KEY COMMANDS
		/* TTL */
		result, err := keyy.Ttl(req)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.IntResponse(result)
	} else if icmd == command.KEYS_I {
		/* KEYS */
		result, err := keyy.Keys(req)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.StrArrResponse(result)
	} else if icmd == command.DEL_I {
		/* DEL */
		rslt, err := keyy.Del(req)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.IntResponse(rslt)
	} else if icmd == command.EXISTS_I {
		/* EXISTS */
		rslt, err := keyy.Exists(req)
		if err != nil {
			return resp.ErrResponse(err), err
		}
		finalRes = resp.StrArrResponse(rslt)
	} else { //Invalid command
		l.Err(errors.New("Invalid cmd:" + string(icmd)))
		return resp.ErrResponse(errors.New(consts.RES_INVALID_CMD)), errors.New(consts.RES_INVALID_CMD)
	}

	return finalRes, nil
}
