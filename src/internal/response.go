/***
** author: duan hongxing
** data: 14 Jul 2024
** desc:
**
**/

package internal

import (
	"encoding/json"
	"topiik/internal/proto"
)

func ErrorResponse(err error) []byte {
	return Response(-1, []byte(err.Error()))
}

func StringResponse(res string, CMD string, msg []byte) (result []byte) {
	buf := []byte(res)
	result = append(result, byte(int8(1)))
	result = append(result, buf...)
	result, _ = proto.Encode(string(result))
	return result
}

func IntegerResponse(res int64, CMD string, msg []byte) (result []byte) {
	buf := byte(res)
	result = append(result, byte(int8(1)))
	result = append(result, buf)
	result, _ = proto.Encode(string(result))
	return result
}

func Response[T any](flag int8, res T) (result []byte) {
	buf, _ := json.Marshal(res)
	result = append(result, byte(flag))
	result = append(result, buf...)
	result, _ = proto.Encode(string(result))
	return result
}
