//author: Duan HongXing
//date: 26 Jun 2024

package persistence

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"topiik/internal/consts"
	"topiik/internal/util"
)

//var lineFeed = byte('\n')

//
// Binary format file, with each msg, 8 Sequence + 4 bytes length + msg
// +-----------------------------8 bytes Sequence---------------------------|--------4 bytes lenght---------------|------msg--------+
// |00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 | 00000000 00000000 00000000 00000000 | xxxxxxxxxxxxxxxx|
// +------------------------------------------------------------------------|-------------------------------------|-----------------+

// active log file
var activeLF *os.File

//func msgId() string {
//	return fmt.Sprintf("%s%v", node.GetNodeInfo().GroupId, util.GetUtcEpoch())
//}

// Append msg to binary log
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
	err = binary.Write(activeLF, binary.LittleEndian, buf)
	if err != nil {
		l.Err(err).Msgf("persist::Append %s", err.Error())
		return err
	}

	// 4: sync to follower(s)
	// TODO: optimize to half nodes in-sync
	syncOne(buf)
	return nil
}

func syncOne(binlogs []byte) {
	/*var err error
	for _, nd := range node.GetPnt().NodeSet {
		var conn *net.TCPConn
		if nd.Id == node.GetNodeInfo().Id { // not sync current node
			continue
		}
		if _, ok := connCache[nd.Id]; ok {
			conn = connCache[nd.Id]
		} else {
			conn, err = util.PreapareSocketClient(nd.Addr2)
			if err != nil {
				l.Err(err).Msgf("persistence::sync conn: %s", err.Error())
				continue
			}
		}
		flrSeq, err := doSync(conn, binlogs)
		if err != nil {
			continue
		}
		if flrSeq < binlogSeq && !batchSyncing {
			go syncBatch(conn, flrSeq+1)
		}
	}*/
}

const batchSyncSize = 1024 * 1024 //

// Load binary log to memory when server start
// Parameters:
//   - execute fn: the func that execute command, i.e: executor.Executer1
func Load(execute execute1) {
	fpath := getActiveBinlogFile()
	exist, err := util.PathExists(fpath)
	if err != nil {
		l.Panic().Msg("[X]load binlog failed")
	}
	if exist {
		f, err := os.OpenFile(fpath, os.O_RDONLY, 0644)
		if err != nil {
			l.Panic().Msg("[X]load binlog failed")
		}
		defer f.Close()

		for {
			buf, err := parseOne(f, &binlogSeq)
			if err != nil {
				if err != io.EOF {
					l.Panic().Msg(err.Error())
				}
				break
			}

			// load to RAM
			replay(buf, execute)
		}
	}
	l.Info().Msgf("persistence::Load BINLOG SEQ: %v", binlogSeq)
}

func getActiveBinlogFile() string {
	return util.GetMainPath() + consts.SLASH + consts.DATA_DIR + consts.SLASH + "000001.bin"
}
