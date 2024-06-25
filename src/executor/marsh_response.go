/***
** author: duan hongxing
** date: 23 Jun 2024
** desc:
**
**/

package executor

import (
	"encoding/json"
	"topiik/internal/datatype"
)

func responseError(err error) []byte {
	return response[string](false, err.Error())
}

func responseSuccess[T any](result T) []byte {
	return response[T](true, result)
}

func response[T any](success bool, response T) []byte {
	b, _ := json.Marshal(&datatype.Response[T]{R: success, M: response})
	return b
}
