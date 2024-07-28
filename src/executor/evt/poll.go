/*
* author: duan hongxing
* date: 28 Jul, 2024
* desc
 */

package evt

import (
	"errors"
	"topiik/internal/datatype"
	"topiik/resp"
)

func poll(req datatype.Req) (res []byte, err error) {
	if len(req.KEYS) != 1 {
		return nil, errors.New(resp.RES_SYNTAX_ERROR)
	}
	return res, nil
}
