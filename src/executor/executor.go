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
	RES_OK           = "OK"
	RES_NIL          = "NIL"
	RES_SYNTAX_ERROR = "ERR:SYNTAX"
	RES_INVALID_OP   = "ERR:INVALID_OP"

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
		pieces, err := needKEY(strs[1])
		if err != nil {
			return marshalResponseError(err)
		}
		result, err := get(pieces)
		if err != nil {
			return marshalResponseError(err)
		}
		return marshalResponseSuccess(result)
	} else if CMD == command.SET {
		/***String GET***/
		pieces, err := needKEY(strs[1])
		if err != nil {
			return marshalResponseError(err)
		}
		result, err := set(pieces)
		if err != nil {
			return marshalResponseError(err)
		}
		return marshalResponseSuccess(result)
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
	return marshalResponseError(errors.New(RES_SYNTAX_ERROR))

}

/***
** Parse the command line after stripped CMD, for commands that KEY is required
** Return:
**	- pieces: [0] the KEY, [1] the request params if have
**	- error: if len(pieces) < 1, means no KEY
**/
func needKEY(keyAndParams string) (pieces []string, err error) {
	pieces = strings.SplitN(strings.TrimLeft(keyAndParams, consts.SPACE), consts.SPACE, 2)
	if len(pieces) < 1 {
		return nil, errors.New(RES_SYNTAX_ERROR)
	}
	return pieces, nil
}

func marshalResponseError(err error) []byte {
	return marshalResponse(err.Error(), false)
}

func marshalResponseSuccess(response string) []byte {
	return marshalResponse(response, true)
}

func marshalResponse(response string, success bool) []byte {
	b, _ := json.Marshal(&datatype.Response{S: success, R: response})
	return b
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
