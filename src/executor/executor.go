/***
* author: duan hongxing
* date: 21 Jun 2024
* desc:
 */

package executor

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"topiik/internal/command"
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/raft"
	"topiik/shared"
)

const (
	WRONG_CMD_MSG = "Wrong command format: "
	INVALID_CMD   = "Invalid command"
)

/***Command RESponse***/
const (
	RES_OK                   = "OK"
	RES_NIL                  = "NIL"
	RES_WRONG_ARG            = "WRONG_ARG"
	RES_WRONG_NUMBER_OF_ARGS = "WRONG_NUM_OF_ARGS"
	RES_DATA_TYPE_NOT_MATCH  = "DATA_TYPE_NOT_MATCH"
	RES_SYNTAX_ERROR         = "SYNTAX_ERR"
	RES_INVALID_CMD          = "INVALID_CMD"
	RES_INVALID_OP           = "INVALID_OP"

	/*** VOTE response ***/
	RES_ACCEPTED = "A"
	RES_REJECTED = "R"
)

func Execute(msg string, serverConfig *config.ServerConfig, nodestatus *raft.NodeStatus) []byte {
	// split into command + arg
	strs := strings.SplitN(strings.TrimLeft(msg, consts.SPACE), consts.SPACE, 2)
	CMD := strings.ToUpper(strings.TrimSpace(strs[0]))
	//result := RES_OK

	if CMD == command.GET { // STRING COMMANDS
		/***String SET***/
		pieces, err := needKEY(strs)
		if err != nil {
			return returnError(err)
		}
		result, err := get(pieces)
		if err != nil {
			return returnError(err)
		}
		return returnSuccessPersistMsg(result, msg)
	} else if CMD == command.SET {
		/***String GET***/
		pieces, err := needKEY(strs)
		if err != nil {
			return returnError(err)
		}
		result, err := set(pieces)
		if err != nil {
			return returnError(err)
		}
		return returnSuccessPersistMsg(result, msg)
	} else if CMD == command.GETM {
		pieces, err := needKEY(strs)
		if err != nil {
			return returnError(err)
		}
		result, err := getM(pieces)
		if err != nil {
			return returnError(err)
		}
		return returnSuccessPersistMsg(result, msg)
	} else if CMD == command.SETM {
		pieces, err := needKEY(strs)
		if err != nil {
			return returnError(err)
		}
		result, err := setM(pieces)
		if err != nil {
			return returnError(err)
		}
		return returnSuccessPersistMsg(result, msg)
	} else if CMD == command.INCR {
		pieces, err := needKEY(strs)
		if err != nil {
			return returnError(err)
		}
		result, err := INCR(pieces)
		if err != nil {
			return returnError(err)
		}
		return returnSuccessPersistMsg(result, msg)
	} else if CMD == command.LPUSH || CMD == command.LPUSHR { // LIST COMMANDS
		/***List LPUSH***/
		pieces, err := needKEY(strs)
		if err != nil {
			return returnError(err)
		}
		result, err := pushList(pieces, CMD)
		if err != nil {
			return returnError(err)
		}
		return returnSuccessPersistMsg(result, msg)
	} else if CMD == command.LPOP || CMD == command.LPOPR {
		pieces, err := needKEY(strs)
		if err != nil {
			return returnError(err)
		}
		result, err := popList(pieces, CMD)
		if err != nil {
			return returnError(err)
		}
		return returnSuccessPersistMsg(result, msg)
	} else if CMD == command.VOTE { // CLUSTER COMMANDS
		if len(strs) != 2 {
			fmt.Printf("%s %s", WRONG_CMD_MSG, msg)
			return []byte(RES_SYNTAX_ERROR)
		} else {
			cTerm, err := strconv.Atoi(strs[1])
			if err != nil {
				return []byte(RES_SYNTAX_ERROR)
			} else {
				return []byte(vote(cTerm, nodestatus))
			}
		}

	} else if CMD == command.APPEND_ENTRY {
		appendEntry(serverConfig, nodestatus)
	} else {
		fmt.Printf("Invalid cmd: %s\n", CMD)
	}
	return returnError(errors.New(RES_SYNTAX_ERROR))

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
func enqueuePersistentMsg(msg string) {
	if shared.MemMap[consts.PERSISTENT_BUF_QUEUE] == nil {
		shared.MemMap[consts.PERSISTENT_BUF_QUEUE] = &datatype.TValue{
			Type:   datatype.TTYPE_LIST,
			TList:  list.New(),
			Expire: consts.UINT32_MAX,
		}
	}
	shared.MemMap[consts.PERSISTENT_BUF_QUEUE].TList.PushFront(msg)
}

/*** Response ***/
func returnError(err error) []byte {
	return response[string](false, err.Error())
}

func returnSuccess[T any](result T) []byte {
	return response[T](true, result)
}

func returnSuccessPersistMsg[T any](result T, msg string) []byte {
	enqueuePersistentMsg(msg)
	return response[T](true, result)
}

func response[T any](success bool, response T) []byte {
	b, _ := json.Marshal(&datatype.Response[T]{R: success, M: response})
	return b
}

