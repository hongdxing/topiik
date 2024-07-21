package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"topiik/internal/command"
)

/*
** Encode string command to byte array
**
 */
func EncodeCmd(strCmd string) ([]byte, error) {
	if icmd, ok := command.CmdCode[strCmd]; ok {
		byteBuf := new(bytes.Buffer)
		binary.Write(byteBuf, binary.LittleEndian, icmd)
		return byteBuf.Bytes(), nil
	}
	return nil, errors.New("invalid command")
}
