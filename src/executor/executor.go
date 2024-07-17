/***
* author: duan hongxing
* date: 21 Jun 2024
* desc:
 */

package executor

import (
	"errors"
	"fmt"
	"strings"
	"topiik/cluster"
	"topiik/internal/command"
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/internal/util"
	"topiik/resp"
)

/***Command RESponse***/
const (
	RES_OK                   = "OK"
	RES_NIL                  = "NIL"
	RES_WRONG_ARG            = "WRONG_ARG"
	RES_WRONG_NUMBER_OF_ARGS = "WRONG_NUM_OF_ARGS"
	RES_DATA_TYPE_NOT_MATCH  = "DATA_TYPE_NOT_MATCH"
	RES_SYNTAX_ERROR         = "SYNTAX_ERR"
	RES_KEY_NOT_EXIST        = "KEY_NOT_EXIST"
	RES_KEY_EXIST_ALREADY    = "KEY_EXIST_ALREADY"

	RES_INVALID_OP = "INVALID_OP"

	RES_INVALID_ADDR = "INVALID_ADDR"

	/*** VOTE response ***/
	RES_ACCEPTED = "A"
	RES_REJECTED = "R"
)

func Execute(msg []byte, srcAddr string, serverConfig *config.ServerConfig) []byte {
	strMsg := msg[4:]
	// split msg into [CMD, params]
	strs := strings.SplitN(strings.TrimLeft(string(strMsg), consts.SPACE), consts.SPACE, 2)
	CMD := strings.ToUpper(strings.TrimSpace(strs[0]))

	// cluster command
	if CMD == command.CLUSTER {
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
	err := srcFilter(srcAddr)
	if err != nil {
		return resp.ErrorResponse(err)
	}

	if CMD == command.GET { // STRING COMMANDS
		/***String SET***/
		//pieces, err := needKEY(strs)
		pieces := []string{}
		if len(strs) == 2 {
			pieces = strings.Split(strs[1], consts.SPACE)
		}
		result, err := get(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result, CMD, msg)
	} else if CMD == command.SET {
		/***String GET***/
		pieces := []string{}
		if len(strs) == 2 {
			pieces = strings.Split(strs[1], consts.SPACE)
		}
		result, err := set(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringResponse(result, CMD, msg)
	} else if CMD == command.GETM {
		//pieces, err := needKEY(strs)
		pieces := []string{}
		if len(strs) == 2 {
			pieces = strings.Split(strs[1], consts.SPACE)
		}
		result, err := getM(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringArrayResponse(result, CMD, msg)
	} else if CMD == command.SETM {
		//pieces, err := needKEY(strs)
		pieces := []string{}
		if len(strs) == 2 {
			pieces = strings.Split(strs[1], consts.SPACE)
		}
		result, err := setM(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.IntegerResponse(int64(result), CMD, msg)
	} else if CMD == command.INCR {
		pieces, err := needKEY(strs)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		result, err := incr(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.IntegerResponse(result, CMD, msg)
	} else if CMD == command.LPUSH || CMD == command.LPUSHR { // LIST COMMANDS
		/***List LPUSH***/
		pieces, err := needKEY(strs)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		result, err := pushList(pieces, CMD)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.IntegerResponse(int64(result), CMD, msg)
	} else if CMD == command.LPOP || CMD == command.LPOPR {
		pieces, err := needKEY(strs)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		result, err := popList(pieces, CMD)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringArrayResponse(result, CMD, msg)
	} else if CMD == command.LLEN {
		pieces, err := needKEY(strs)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		result, err := llen(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.IntegerResponse(int64(result), CMD, msg)
	} else if CMD == command.TTL { // KEY COMMANDS
		pieces := splitParams(strs)
		result, err := ttl(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.IntegerResponse(result, CMD, msg)
	} else if CMD == command.KEYS {
		pieces := splitParams(strs)
		result, err := keys(pieces)
		if err != nil {
			return resp.ErrorResponse(err)
		}
		return resp.StringArrayResponse(result, CMD, msg)
	} else {
		fmt.Printf("Invalid cmd: %s\n", CMD)
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
** Split command parameters if any
**
** Return:
**	The command pieces except the CMD itself
**/
func splitParams(strs []string) (pieces []string) {
	if len(strs) == 2 {
		pieces = strings.Split(strs[1], consts.SPACE)
	}
	return pieces
}

/***
** Persist command
**
**/
/*
func enqueuePersistentMsg(msg []byte) {
	if shared.MemMap[consts.PERSISTENT_BUF_QUEUE] == nil {
		shared.MemMap[consts.PERSISTENT_BUF_QUEUE] = &datatype.TValue{
			Typ: datatype.V_TYPE_LIST,
			Lst: list.New(),
			Exp: consts.UINT32_MAX,
		}
	}
	shared.MemMap[consts.PERSISTENT_BUF_QUEUE].Lst.PushFront(msg)
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