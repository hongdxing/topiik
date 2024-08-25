/*
* author: Duan Hongxing
* date: 22 Aug, 2024
* desc:
*
 */
package clus

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
func NewPartition(req datatype.Req) (ptnIds []string, err error) {
	if len(req.Args) == 0 {
		return nil, errors.New(resp.RES_SYNTAX_ERROR)
	}
	var ptnCount int
	ptnCount, err = strconv.Atoi(req.Args)
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
