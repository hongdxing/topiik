// ©2024 www.topiik.com
// author: Duan Hongxing
// date: 1 Sep, 2024

package server

import (
	"errors"
	"topiik/cluster"
	"topiik/resp"
)

// Check if new online node in the cluster or not
func online(pieces []string) (string, error) {

	ndId := pieces[0]
	if _, ok := cluster.GetControllerInfo().Nodes[ndId]; ok {
		return resp.RES_OK, nil
	}
	if _, ok := cluster.GetControllerInfo().Nodes[ndId]; ok {
		return resp.RES_OK, nil
	}

	return resp.RES_REJECTED, errors.New(resp.RES_REJECTED)
}
