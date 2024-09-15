/*
* ©2024 www.topiik.com
* Author: Duan HongXing
* Date: 30 Aug, 2024
 */

package clus

import (
	"errors"
	"strings"
	"topiik/internal/datatype"
	"topiik/resp"
)

/*
* Remove node, Controller or Worker from cluster
* Syntax: REMOVE-NODE nodeId
 */
func RemoveNode(req datatype.Req) error {
	ndId := strings.TrimSpace(req.Args)
	if len(ndId) == 0 {
		return errors.New(resp.RES_SYNTAX_ERROR)
	}
	return nil
	//return cluster.RemoveNode(ndId)
}
