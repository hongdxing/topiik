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
func ReceiveBinlog(binlogs []byte, f execute1) (seq int64, err error) {
	if flrActiveBLF == nil {
		path := getActiveBinlogFile()
		flrActiveBLF, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			l.Err(err).Msgf("persistence::ReceiveBinlog %s", err.Error())
		}
	}
	//var seq int64
	reader := bytes.NewReader(binlogs)
	for {
		buf, err := parseOne(reader, &seq)
		if err != nil {
			if err != io.EOF {
				l.Err(err).Msg(err.Error())
			}
			break
		}
		if (seq - 1) != binlogSeq {
			l.Warn().Msgf("persistence::ReceiveBinlog binlog lag %v", seq-binlogSeq)
			break
		}

		err = binary.Write(flrActiveBLF, binary.LittleEndian, buf)
		if err != nil {
			seq-- // parse ok, but write failed, return pre seq
			break
		}

		/* load to RAM */
		replay(buf, f)

		/* update global seq */
		binlogSeq = seq
	}
	return binlogSeq, nil
}
