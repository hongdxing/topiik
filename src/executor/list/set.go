// author: Duan HongXing
// date: 4 Sep, 2024

package list

import (
	"errors"
	"strconv"
	"strings"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

// Set value at index of list
// Syntax: LSET Key Value Index
func Set(req datatype.Req) (rslt string, err error) {
	key := string(req.Keys[0])
	//key := bytes.NewBuffer(req.Keys[0]).String()
	if val, ok := memo.MemMap[key]; ok {
		// return RES_DATA_TYPE_MISMATCH if the key is not List
		if val.Typ != memo.V_TYPE_LIST {
			return "", errors.New(resp.RES_DATA_TYPE_MISMATCH)
		}

		//
		pieces := strings.Split(req.Args, consts.SPACE)
		if len(pieces) != 1 {
			return "", errors.New(resp.RES_SYNTAX_ERROR)
		}

		// Parse index
		idx, err := strconv.Atoi(pieces[0])
		if err != nil {
			return "", errors.New(resp.RES_SYNTAX_ERROR)
		}

		// Convert minus
		// e.g. len=5, idx=-1 => idx=-1+5=4
		if idx < 0 {
			idx = idx + val.Lst.Len()
		}

		// Out of bound
		if idx >= val.Lst.Len() || idx < 0 {
			return "", errors.New(resp.RES_OUT_OF_BOUND)
		}

		var i int
		for ele := val.Lst.Front(); ele != nil; ele = ele.Next() {
			// Double check
			if i >= val.Lst.Len() {
				return "", errors.New(resp.RES_KEY_NOT_EXIST)
			}
			if i < idx {
				i++
				continue
			} else if i == idx {
				ele.Value = req.Vals[0]
			}
		}

	} else {
		return "", errors.New(resp.RES_KEY_NOT_EXIST)
	}

	return resp.RES_OK, nil
}
