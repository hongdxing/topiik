/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package cluster

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"topiik/internal/command"
	"topiik/internal/config"
	"topiik/internal/consts"
)

const (
	RES_SYNTAX_ERROR = "SYNTAX_ERR"
)

func Execute(msg []byte, serverConfig *config.ServerConfig) (result []byte, err error) {
	strs := strings.SplitN(strings.TrimLeft(string(msg[4:]), consts.SPACE), consts.SPACE, 2)
	CMD := strings.ToUpper(strings.TrimSpace(strs[0]))

	if CMD == command.CLUSTER_JOIN_ACK {
		pieces := splitParams(strs)
		if len(pieces) < 1 {
			return nil, errors.New(RES_SYNTAX_ERROR)
		}
		result, err := clusterJoin(pieces)
		if err != nil {
			return nil, err
		}
		return []byte(result), nil
	} else if CMD == "VOTE" {
		if len(strs) != 2 {
			fmt.Printf("%s %s", RES_SYNTAX_ERROR, msg)
			return nil, errors.New(RES_SYNTAX_ERROR)
		} else {
			cTerm, err := strconv.Atoi(strs[1])
			if err != nil {
				return nil, errors.New(RES_SYNTAX_ERROR)
			} else {
				return []byte(vote(cTerm)), nil
			}
		}
	} else if CMD == RPC_APPENDENTRY {
		pieces := splitParams(strs)
		err := appendEntry(pieces, serverConfig)
		if err != nil {
			return nil, err
		}
		return []byte{}, nil
	} else {
		// forward msg to Workers
		Forward(msg)
	}
	return nil, errors.New(RES_SYNTAX_ERROR)
}

func splitParams(strs []string) (pieces []string) {
	if len(strs) == 2 {
		pieces = strings.Split(strs[1], consts.SPACE)
	}
	return pieces
}
