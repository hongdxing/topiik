//author: Duan Hongxing
//date: 22 Aug, 2024

package clus

import (
	"errors"
	"strings"
	"topiik/cluster"
	"topiik/internal/datatype"
	"topiik/resp"
)

// Create new partition
// Only after RESHARD the new partition  before
// Syntax: NEW-PARTITION host:port[,host:port...]
func NewPartition(req datatype.Req) (ptnId string, err error) {
	if len(req.Args) == 0 {
		return "", errors.New(resp.RES_SYNTAX_ERROR)
	}

	workers := strings.Split(req.Args, ",")

	wrkNodeIdAddr, _ := checkConnection(workers)
	ptnId, err = cluster.NewPartition(wrkNodeIdAddr)
	if err != nil {
		return "", err
	}

	return ptnId, nil
}
