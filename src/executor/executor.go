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
	RES_ERROR_CMD    = "ErrCMD"
	RES_SYNTAX_ERROR = "ERR:SYNTAX"

	/*** VOTE response ***/
	RES_ACCEPTED = "A"
	RES_REJECTED = "R"
)

var memMap = make(map[string]*datatype.TValue)

func Execute(msg string, serverConfig *config.ServerConfig, nodestatus *raft.NodeStatus) string {
	// split into command + arg
	strs := strings.SplitN(msg, " ", 2)
	CMD := strings.TrimSpace(strs[0])
	result := RES_OK

	if CMD == command.GET {
		/***String SET***/
		set(strs[1])
	} else if CMD == command.SET {
		/***String GET***/
	} else if CMD == command.VOTE {
		//conn.Write([]byte(command.RES_REJECTED))
		if len(strs) != 2 {
			fmt.Printf("%s %s", WRONG_CMD_MSG, msg)
			result = RES_ERROR_CMD
		} else {
			cTerm, err := strconv.Atoi(strs[1])
			if err != nil {
				result = RES_ERROR_CMD
			} else {
				result = vote(cTerm, nodestatus)
				//conn.Write([]byte(result))
			}
		}

	} else if CMD == command.APPEND_ENTRY {
		appendEntry(serverConfig, nodestatus)
	} else {
		fmt.Printf("Invalid cmd: %s\n", CMD)
	}
	return result

}

func parseSingleValueCMD(strs []string) (string, error) {
	if len(strs) != 2 {
		fmt.Printf("%s", WRONG_CMD_MSG)
		return "", errors.New(INVALID_CMD)
	}
	return strs[1], nil
}
