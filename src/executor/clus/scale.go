/*
* author: duan hongxing
* date: 25 Jul 2024
* desc:
 */

package clus

import (
	"errors"
	"strconv"
	"strings"
	"topiik/cluster"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/resp"
)

func Scale(req datatype.Req) (result string, err error) {
	l.Info().Msg("executor::scale start")
	var partition = 0
	var replica = 0
	pieces := strings.Split(req.Args, consts.SPACE)
	for i := 0; i < len(pieces); i++ {
		if strings.ToLower(pieces[i]) == "partition" {
			if len(pieces) > i {
				partition, err = strconv.Atoi(pieces[i+1])
				if err != nil {
					return "", errors.New(resp.RES_SYNTAX_ERROR)
				}
			}
			i++
		} else if strings.ToLower(pieces[i]) == "replica" {
			if len(pieces) > i {
				replica, err = strconv.Atoi(pieces[i+1])
				if err != nil {
					return "", errors.New(resp.RES_SYNTAX_ERROR)
				}
			}
			i++
		} else {
			return "", errors.New(resp.RES_SYNTAX_ERROR)
		}
	}
	if partition <= 0 || replica <= 0 {
		return "", errors.New(resp.RES_SYNTAX_ERROR)
	}

	result, err = cluster.Scale(partition, replica)
	if err != nil {
		return "", err
	}

	l.Info().Msg("executor::scale done")
	return result, nil
}
