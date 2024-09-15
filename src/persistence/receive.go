//author: Duan Hongxing
//date: 13 Aug, 2024

package persistence

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

// follower active binlog file
var flrActiveBLF *os.File

// Parameters:
//   - binlogs: one or more logs
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
		//buf, err := parseOne(reader, &seq)
		// binlog
		bl, err := receiveOne(reader)
		if err != nil {
			if err != io.EOF {
				l.Err(err).Msg(err.Error())
			}
			break
		}
		// prepend seq
		binlogSeq += 1
		bbuf := new(bytes.Buffer)
		binary.Write(bbuf, binary.LittleEndian, binlogSeq)
		buf := bbuf.Bytes()
		buf = append(buf, bl...)

		err = binary.Write(flrActiveBLF, binary.LittleEndian, buf)
		if err != nil {
			// ??? need deduct???
			//binlogSeq-- // parse ok, but write failed, return pre seq
			break
		}

		/* load to RAM */
		//replay(buf, f)

		/* update global seq */
		//binlogSeq = seq
	}
	return binlogSeq, nil
}

func receiveOne(r io.Reader) (res []byte, err error) {
	var length int32

	// 1: read length
	buf := make([]byte, 4)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msg(err.Error())
		}
		return nil, err
	}
	bbuf := bytes.NewBuffer(buf)
	binary.Read(bbuf, binary.LittleEndian, &length)
	res = append(res, buf...)

	// 2: read msg
	buf = make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msg(err.Error())
		}
		return nil, err
	}
	res = append(res, buf...)
	return res, nil
}
