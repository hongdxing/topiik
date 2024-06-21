package executor

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"topiik/internal/command"
	"topiik/raft"
)

const (
	WRONG_CMD_MSG = "Wrong command format: "
)

const (
	RES_OK        = "OK"
	RES_ERROR_CMD = "ErrCMD"

	/*** VOTE response ***/
	RES_ACCEPTED = "A"
	RES_REJECTED = "R"
)

func Execute(conn net.Conn, msg string, nodestatus *raft.NodeStatus) {
	// split into command + arg
	strs := strings.SplitN(msg, " ", 2)
	CMD := strings.TrimSpace(strs[0])
	result := RES_OK
	if CMD == command.VOTE {
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
				conn.Write([]byte(result))
			}
		}
	}

}
