/*
* author: duan hongxing
* data: 14 Jul 2024
* desc:
*	for return response to client
* 	there are 3 response types:
*	- 1: string
*	- 2: integer
*	- 3: string array
* 	the format is:
*	------------Reposne Lenght(4)---------Flag(1)---Datatype(1)---Body---
*	00000000 00000000 00000000 00000000 | 00000000 | 00000000 |   ...
*	---------------------------------------------------------------------
*
 */

package resp

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"topiik/internal/proto"
)

func ErrResponse(err error) (result []byte) {
	result = append(result, byte(int8(0)))
	result = append(result, byte(int8(1))) // string type
	result = append(result, []byte(err.Error())...)
	result, _ = proto.Encode(string(result))
	return result
}

func StrResponse(res string) (result []byte) {
	buf := []byte(res)
	result = append(result, byte(int8(1)))
	result = append(result, byte(int8(1))) // 1: string type
	result = append(result, buf...)
	result, _ = proto.Encode(string(result))
	return result
}

func IntResponse(res int64) (result []byte) {
	var buffer = new(bytes.Buffer)
	// Write message HEADER
	err := binary.Write(buffer, binary.LittleEndian, int8(1)) // one byte of success flag
	if err != nil {
		l.Err(err).Msgf("IntResponse() write flag:%s", err.Error())
		return ErrResponse(err)
	}
	err = binary.Write(buffer, binary.LittleEndian, int8(2)) // 2: integer type
	if err != nil {
		l.Err(err).Msgf("IntResponse() write type:%s", err.Error())
		return ErrResponse(err)
	}

	err = binary.Write(buffer, binary.LittleEndian, res)
	if err != nil {
		l.Err(err).Msgf("IntResponse() write res:%s", err.Error())
		return ErrResponse(err)
	}
	result, err = proto.Encode(buffer.String())
	if err != nil {
		return ErrResponse(err)
	}
	return result
}

func StrArrResponse(res []string) (result []byte) {
	buf, _ := json.Marshal(res)
	result = append(result, byte(int8(1)))
	result = append(result, byte(int8(3))) // 3: string array type
	result = append(result, buf...)
	result, _ = proto.EncodeB(result)
	return result
}

func RedirectResponse(leaderAddr string) (result []byte) {
	result = append(result, byte(int8(1)))
	result = append(result, byte(int8(32))) // why 32? 302 = http redirect, but int8 not enough ~!~
	result = append(result, []byte(leaderAddr)...)
	result, _ = proto.Encode(string(result))
	return result
}

/*
func Response[T any](flag int8, res T) (result []byte) {
	buf, _ := json.Marshal(res)
	result = append(result, byte(flag))
	result = append(result, buf...)
	result, _ = proto.Encode(string(result))
	return result
}
*/

/*
* Return success/fail flag of response
*
 */
func ParseResFlag(res []byte) int8 {
	flagByte := res[4:5]
	bbuf := bytes.NewBuffer(flagByte)
	var flag int8
	err := binary.Read(bbuf, binary.LittleEndian, &flag)
	if err != nil {
		return 0
	}
	return flag
}
