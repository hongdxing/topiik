/*
* @author: Duan Hongxing
* @date: 22 Aug, 2024
* @desc:
*
 */
package executor

import (
	"errors"
	"strconv"
	"topiik/cluster"
	"topiik/internal/datatype"
	"topiik/resp"
)

/*
* Create new partition
* Syntax: NEW-PARTITION count
 */
func newPartition(req datatype.Req) (ptnIds []string, err error) {
	if len(req.ARGS) == 0 {
		return nil, errors.New(resp.RES_SYNTAX_ERROR)
	}
	var ptnCount int
	ptnCount, err = strconv.Atoi(req.ARGS)
	if err != nil {
		return nil, err
	}
	if ptnCount <= 0 {
		return nil, errors.New(resp.RES_SYNTAX_ERROR)
	}
	ptnIds, err = cluster.NewPartition(ptnCount)
	if err != nil {
		return nil, err
	}

	return ptnIds, nil
}
