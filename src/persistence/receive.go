/*
* author: duan hongxing
* date: 13 Aug, 2024
* desc:
*	Receive binlog from partition leader
 */

package persistence

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

// follower active binlog file
var flrActiveBLF *os.File

/*
* Parameters:
*	- binlogs: one or more logs
*
 */
func ReceiveBinlog(binlogs []byte) (res []byte, err error) {
	if flrActiveBLF == nil {
		path := getActiveBinlogFile()
		flrActiveBLF, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			l.Err(err).Msgf("persistence::ReceiveBinlog %s", err.Error())
		}
	}
	var seq int64
	reader := bytes.NewReader(binlogs)
	for {
		buf, err := parseOne(reader, &seq)
		if err != nil {
			if err != io.EOF {
				l.Err(err).Msg(err.Error())
			}
			break
		}
		err = binary.Write(flrActiveBLF, binary.LittleEndian, buf)
		if err != nil {
			seq-- // parse ok, but write failed, return pre seq
			break
		}
	}
	byteBuf := new(bytes.Buffer)
	binary.Write(byteBuf, binary.LittleEndian, seq)
	res = byteBuf.Bytes()
	return res, nil
}
