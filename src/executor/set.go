/***
* author: duan hongxing
* date: 22 Jun 2024
* desc:
***/

package executor

import (
	"strings"
	"topiik/internal/consts"
	"topiik/internal/datatype"
)

func set(params string) string {
	strs := strings.SplitN(params, consts.SPACE, 2)
	if len(strs) == 2 {
		memMap[strings.TrimSpace(strs[0])] = &datatype.TValue{
			Type:   datatype.TTYPE_STRING,
			String: []byte(strs[1]),
			Expire: consts.UINT32_MAX}
		return RES_OK
	} else {
		return RES_SYNTAX_ERROR
	}

}
