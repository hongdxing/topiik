/*
* author: duan hongxing
* date: 25 Jul 2024
* desc:
 */

package executor

import (
	"errors"
	"strconv"
	"strings"
	"topiik/cluster"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/resp"
)

func scale(req datatype.Req) (result string, err error) {
	log.Info().Msg("***scale start***")
	pieces := strings.SplitN(req.ARGS, consts.SPACE, 2)
	if len(pieces) != 2 {
		return "", errors.New(resp.RES_SYNTAX_ERROR)
	}
	partitions, err := strconv.ParseUint(pieces[0], 10, 16)
	if err != nil {
		return "", errors.New(resp.RES_SYNTAX_ERROR)
	}
	replicas, err := strconv.ParseUint(pieces[1], 10, 16)
	if err != nil {
		return "", errors.New(RES_SYNTAX_ERROR)
	}

	cluster.Scale(int(partitions), int(replicas))

	log.Info().Msg("***scale done***")
	return result, nil
}
