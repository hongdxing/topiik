/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package ccss

import (
	"errors"
	"fmt"
	"strings"
	"topiik/internal/command"
	"topiik/internal/consts"
)

const (
	RES_SYNTAX_ERROR = "SYNTAX_ERR"
)

func Execute(msg []byte) (result []byte, err error) {
	strs := strings.SplitN(strings.TrimLeft(string(msg[4:]), consts.SPACE), consts.SPACE, 2)
	CMD := strings.ToUpper(strings.TrimSpace(strs[0]))

	if CMD == "CLUSTER" {
		pieces := splitParams(strs)
		if len(pieces) < 1 {
			return nil, errors.New(RES_SYNTAX_ERROR)
		}
		fmt.Println(pieces)
		if strings.ToUpper(pieces[0]) == "INFO" {
			//TODO
		} else if strings.ToUpper(pieces[0]) == command.CLUSTER_JOIN_ACK {
			fmt.Println("---Join ack---")
			result, err := clusterJoin(pieces)
			if err != nil {
				return nil, err
			}
			return []byte(result), nil
		}
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
