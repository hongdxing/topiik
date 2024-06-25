/***
* author: duan hongxing
* date: 21 Jun 2024
* desc:
 */

package executor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"topiik/internal/command"
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/raft"
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

var memMap = make(map[string]*datatype.TValue)

func Execute(msg string, serverConfig *config.ServerConfig, nodestatus *raft.NodeStatus) []byte {
	// split into command + arg
	strs := strings.SplitN(strings.TrimLeft(msg, consts.SPACE), consts.SPACE, 2)
	CMD := strings.ToUpper(strings.TrimSpace(strs[0]))
	//result := RES_OK

	if CMD == command.GET {
		/***String SET***/
		pieces, err := needKEY(strs)
		if err != nil {
			return responseError(err)
		}
		result, err := get(pieces)
		if err != nil {
			return responseError(err)
		}
		return responseSuccess(result)
	} else if CMD == command.SET {
		/***String GET***/
		pieces, err := needKEY(strs)
		if err != nil {
			return responseError(err)
		}
		result, err := set(pieces)
		if err != nil {
			return responseError(err)
		}
		return responseSuccess(result)
	} else if CMD == command.LPUSH || CMD == command.RPUSH {
		/***List LPUSH***/
		pieces, err := needKEY(strs)
		if err != nil {
			return responseError(err)
		}
		result, err := pushList(pieces, CMD)
		if err != nil {
			return responseError(err)
		}
		return responseSuccess(result)

	} else if CMD == command.LPOP || CMD == command.RPOP {
		pieces, err := needKEY(strs)
		if err != nil {
			return responseError(err)
		}
		result, err := popList(pieces, CMD)
		if err != nil {
			return responseError(err)
		}
		return responseSuccess(result)
	} else if CMD == command.VOTE {
		//conn.Write([]byte(command.RES_REJECTED))
		if len(strs) != 2 {
			fmt.Printf("%s %s", WRONG_CMD_MSG, msg)
			return []byte(RES_SYNTAX_ERROR)
		} else {
			cTerm, err := strconv.Atoi(strs[1])
			if err != nil {
				return []byte(RES_SYNTAX_ERROR)
			} else {
				return []byte(vote(cTerm, nodestatus))
				//conn.Write([]byte(result))
			}
		}

	} else if CMD == command.APPEND_ENTRY {
		appendEntry(serverConfig, nodestatus)
	} else {
		fmt.Printf("Invalid cmd: %s\n", CMD)
	}
	return responseError(errors.New(RES_SYNTAX_ERROR))

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

/*
func marshalResponse(response string, err error) []byte {
	success := true
	if err != nil {
		success = false
	}
	b, _ := json.Marshal(&datatype.Response{S: success, R: response})
	return b
}*/
