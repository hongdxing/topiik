/***
** author: duan hongxing
** date: 26 Jun 2024
** desc:
**
**/

package persistence

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
	"topiik/executor"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/logger"
)

var l = logger.Get()

//var lineFeed = byte('\n')

var msgSeq int64 = 0

/*
* Binary format file, with each msg, 8 Sequence + msg
*
*
 */
func Persist() {
	filePath := getCurrentLogFile()
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		l.Panic().Msg(err.Error())
	}
	defer file.Close()

	for {
		// msg
		msg := <-executor.PersistenceCh
		// msg sequence
		msgSeq++

		byteBuf := new(bytes.Buffer)
		// 1: sequence
		byteBuf.Reset()
		binary.Write(byteBuf, binary.LittleEndian, msgSeq)
		buf := byteBuf.Bytes()

		// 2: msg
		buf = append(buf, msg...)

		// write to file
		binary.Write(file, binary.NativeEndian, buf)
	}
}

func Load() {
	filePath := getCurrentLogFile()
	exist, err := util.PathExists(filePath)
	if err != nil {
		l.Panic().Msg("[X]load binlog failed")
	}
	if !exist {
		return
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		l.Panic().Msg("[X]load binlog failed")
	}
	defer file.Close()

	var length int32
	for {
		// 1: read sequence
		buf := make([]byte, 8)
		_, err = io.ReadFull(file, buf)
		if err != nil && err != io.EOF {
			l.Panic().Msg(err.Error())
		} else if err == io.EOF {
			return
		}

		byteBuf := bytes.NewBuffer(buf)
		binary.Read(byteBuf, binary.LittleEndian, &msgSeq)
		l.Info().Msgf("sequence: %v", msgSeq)

		// 2: read lenght
		buf = make([]byte, 4)
		_, err = io.ReadFull(file, buf)
		if err != nil && err != io.EOF {
			l.Panic().Msg(err.Error())
		}

		byteBuf = bytes.NewBuffer(buf)
		binary.Read(byteBuf, binary.LittleEndian, &length)

		// 3: read msg
		buf = make([]byte, length)
		_, err = io.ReadFull(file, buf)
		if err != nil && err != io.EOF {
			l.Panic().Msg(err.Error())
		}

		// 4: execute msg
		icmd, _, err := proto.DecodeHeader(buf)
		if err != nil {
			l.Panic().Msgf("persist::Load %s", err.Error())
		}

		var req datatype.Req
		err = json.Unmarshal(buf[2:], &req) // 2= 1 icmd and 1 ver
		if err != nil {
			l.Panic().Msgf("persist::Load %s", err.Error())
		}
		executor.Execute1(icmd, req)

		if err != nil {
			return
		}
	}
}

/* One line for each command having problem of char LF can break line unintentionaly
func Persist() {
	filePath := getCurrentLogFile()
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		l.Panic().Msg(err.Error())
	}
	defer file.Close()

	for {
		// msg
		msgBytes := <-executor.PersistenceCh
		msgBytes = append(msgBytes, lineFeed) // append line break in the end
		// msg sequence
		msgSeq++
		byteBuf := new(bytes.Buffer)
		binary.Write(byteBuf, binary.BigEndian, msgSeq)
		buf := byteBuf.Bytes()
		buf = append(buf, msgBytes...)
		// write to file
		binary.Write(file, binary.NativeEndian, buf)
	}
}

func Load() {
	filePath := getCurrentLogFile()
	exist, err := util.PathExists(filePath)
	if err != nil {
		l.Panic().Msg("[X]load binlog failed")
	}
	if !exist {
		return
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		l.Panic().Msg("[X]load binlog failed")
	}
	scanner := bufio.NewScanner(file)
	// resize scanner's capacity for lines over 64K, see next example
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		msg := scanner.Bytes()

		if len(msg) == 0 || (len(msg) == 1 && msg[0] == lineFeed) {
			l.Warn().Msgf("persist::Load Empty at %v", msgSeq+1)
			continue
		}

		// validate minium lenght
		if len(msg) < 12 { // 8 bytes of seq, 4 bytes of length header
			l.Panic().Msgf("persist::Load invalid msg %s", msg)
		}

		// get msg sequence
		byteBuf := bytes.NewBuffer(msg[0:8])
		binary.Read(byteBuf, binary.BigEndian, &msgSeq)
		l.Info().Msgf("sequence: %v", msgSeq)

		// remove msg seq
		msg = msg[8:]

		// remove line feed
		last := msg[len(msg)-1]
		if last == lineFeed { // remove last '\n'
			msg = msg[:len(msg)-1]
		}

		// remove msg length header, the final msg to Execute
		msg = msg[4:]

		icmd, _, err := proto.DecodeHeader(msg)
		if err != nil {
			l.Panic().Msgf("persist::Load %s", err.Error())
		}

		var req datatype.Req
		err = json.Unmarshal(msg[2:], &req) // 2= 1 icmd and 1 ver
		if err != nil {
			l.Panic().Msgf("persist::Load %s", err.Error())
		}
		executor.Execute1(icmd, req)
	}

	if err := scanner.Err(); err != nil {
		l.Panic().Msg(err.Error())
	}
}
*/

func getCurrentLogFile() string {
	return util.GetMainPath() + consts.SLASH + consts.DATA_DIR + consts.SLASH + "000001.bin"
}
