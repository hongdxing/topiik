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
	"net"
	"os"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
)

//var lineFeed = byte('\n')

/*
* Binary format file, with each msg, 8 Sequence + msg
* |-----------------------------8 bytes Sequence---------------------------|--------4 bytes lenght---------------|------msg--------|
* |00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 | 00000000 00000000 00000000 00000000 | xxxxxxxxxxxxxxxx|
* |------------------------------------------------------------------------|-------------------------------------|-----------------|
*
 */
/*
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
		binLogSeq++

		// 1: sequence
		byteBuf := new(bytes.Buffer)
		binary.Write(byteBuf, binary.LittleEndian, binLogSeq)
		buf := byteBuf.Bytes()

		// 2: msg
		buf = append(buf, msg...)

		// write to file
		binary.Write(file, binary.NativeEndian, buf)
	}
}
*/

// active log file
var activeLF *os.File

var connCache = make(map[string]*net.TCPConn)

/*
* Append msg to binary log
*
 */
func Append(msg []byte) (err error) {
	if activeLF == nil {
		var err error
		filePath := getActiveBinlogFile()
		activeLF, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			l.Panic().Msg(err.Error())
		}
	}
	binlogSeq++

	// 1: sequence
	bbuf := new(bytes.Buffer)
	binary.Write(bbuf, binary.LittleEndian, binlogSeq)
	buf := bbuf.Bytes()

	// 2: msg
	buf = append(buf, msg...)

	// 3: write to file
	err = binary.Write(activeLF, binary.NativeEndian, buf)
	if err != nil {
		l.Err(err).Msgf("persist::Append %s", err.Error())
		return err
	}

	// 4: sync to follower
	sync(buf)
	return nil
}

/*
* Sync to follower(s)
*
 */
func sync(logitem []byte) (err error) {
	var conn *net.TCPConn
	for _, nd := range node.GetPnt().NodeSet {
		if nd.Id == node.GetNodeInfo().Id { // not sync current node
			continue
		}
		if _, ok := connCache[nd.Id]; ok {
			conn = connCache[nd.Id]
		} else {
			conn, err = util.PreapareSocketClient(nd.Addr2)
			if err != nil {
				l.Err(err).Msgf("persistence::sync %s", err.Error())
			}
		}
		// icmd
		bbuf := new(bytes.Buffer)
		binary.Write(bbuf, binary.LittleEndian, consts.RPC_SYNC_BINLOG)

		// binlog
		buf := bbuf.Bytes()
		buf = append(buf, logitem...)

		// encode
		msg, err := proto.EncodeB(buf)
		if err != nil {
			l.Err(err).Msgf("persistence::sync %s", err.Error())
			continue
		}

		// send
		_, err = conn.Write(msg)
		if err != nil {
			l.Err(err).Msgf("persistence::sync %s", err.Error())
			continue
		}
	}
	return nil
}

type fn func(uint8, datatype.Req) []byte

/*
* Load binary log to memory when server start
* Parameters:
*	- f fn: the func that execute command, i.e: executor.Executer1
 */
func Load(f fn) {
	filePath := getActiveBinlogFile()
	exist, err := util.PathExists(filePath)
	if err != nil {
		l.Panic().Msg("[X]load binlog failed")
	}
	if exist {
		file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
		if err != nil {
			l.Panic().Msg("[X]load binlog failed")
		}
		defer file.Close()

		for {
			buf, err := parseOne(file, &binlogSeq)
			if err != nil {
				if err != io.EOF {
					l.Panic().Msg(err.Error())
				}
				break
			}

			// 4: replay msg(load from disk to memory)
			buf = buf[12:]
			icmd, _, err := proto.DecodeHeader(buf)
			if err != nil {
				l.Panic().Msgf("persistence::Load %s", err.Error())
			}

			var req datatype.Req
			buf = buf[2:]
			err = json.Unmarshal(buf, &req) // 2= 1 icmd and 1 ver
			if err != nil {
				l.Panic().Msgf("persistence::Load %s", err.Error())
			}

			// replay to load to RAM
			f(icmd, req)
		}
	}
	l.Info().Msgf("persistence::Load BINLOG SEQ: %v", binlogSeq)
}

func getActiveBinlogFile() string {
	return util.GetMainPath() + consts.SLASH + consts.DATA_DIR + consts.SLASH + "000001.bin"
}
