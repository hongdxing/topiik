/***
** author: duan hongxing
** data: 14 Jul 2024
** desc:
**
**/

package resp

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"topiik/internal/proto"
)

func ErrorResponse(err error) (result []byte) {
	result = append(result, byte(int8(0)))
	result = append(result, []byte(err.Error())...)
	result, _ = proto.Encode(string(result))
	return result
}

func StringResponse(res string, CMD string, msg []byte) (result []byte) {
	buf := []byte(res)
	result = append(result, byte(int8(1)))
	result = append(result, buf...)
	result, _ = proto.Encode(string(result))
	return result
}

func IntegerResponse(res int64, CMD string, msg []byte) (result []byte) {
	var buffer = new(bytes.Buffer)
	// Write message HEADER
	err := binary.Write(buffer, binary.LittleEndian, int8(1)) // one byte of success flag
	if err != nil {
		fmt.Printf("IntegerResponse() write flag:%s", err.Error())
		return ErrorResponse(err)
	}
	err = binary.Write(buffer, binary.LittleEndian, res)
	if err != nil {
		fmt.Printf("IntegerResponse() write res:%s", err.Error())
		return ErrorResponse(err)
	}
	result, err = proto.Encode(buffer.String())
	if err != nil {
		return ErrorResponse(err)
	}
	return result
}

func StringArrayResponse(res []string, CMD string, msg []byte) (result []byte) {
	buf, _ := json.Marshal(res)
	result = append(result, byte(int8(1)))
	result = append(result, buf...)
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
