//author: Duan HongXing
//date: 26 Jun 2024

package persistence

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/binary"
	"io"
	"net"
	"os"
	"sync"
	"time"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

//var lineFeed = byte('\n')

//
// Binary format file, with each msg, 8 Sequence + 4 bytes length + msg
// +-----------------------------8 bytes Sequence---------------------------|--------4 bytes lenght---------------|------msg--------+
// |00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 | 00000000 00000000 00000000 00000000 | xxxxxxxxxxxxxxxx|
// +------------------------------------------------------------------------|-------------------------------------|-----------------+

// active log file
var activeLF *os.File

// cache tcp conn to follower
var connCache = make(map[string]*net.TCPConn)
var pstConn *net.TCPConn

// is batch syncing in progress
var batchSyncing = false

var msgQ list.List
var lock sync.Mutex
var dequeueTicker time.Ticker

// enqueue msg pending persist to persistor server
func Enqueue(msg []byte) {
	lock.Lock()
	defer lock.Unlock()
	msgQ.PushFront(msg)
}

func Dequeue() {
	dequeueTicker = *time.NewTicker(time.Millisecond * 1000)
	for {
		<-dequeueTicker.C
		doDequeue()
	}
}

func doDequeue() {
	lock.Lock()
	defer lock.Unlock()

	var binlogs []byte
	for {
		if msgQ.Len() == 0 {
			break
		}
		buf := msgQ.Back().Value.([]byte)
		msgQ.Remove(msgQ.Back())
		binlogs = append(binlogs, buf...)
		if len(binlogs) > batchSyncSize { // limit batch size
			break
		}
	}

	// todo: if connection failed???
	if pstConn == nil {
		addr := node.GetConfig().Persistors[0]
		host, _, port2, _ := util.SplitAddress2(addr)
		conn, err := util.PreapareSocketClient(host + ":" + port2)
		pstConn = conn
		if err != nil {
			l.Err(err).Msgf("persistence::doDequeue conn: %s", err.Error())
		}
	}
	if pstConn == nil {
		l.Warn().Msg("persist::doDequeue No PERSISTOR available!!!")
		return
	}

	// if connection closed by remote persistor, then set pstConn to nil,
	// so that next time can re-connect
	_, err := doSync(pstConn, binlogs)
	if err != nil && err == net.ErrClosed {
		pstConn = nil
	}
}

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

// if follower is fall behind, then try to send binlog in batch to folower(s)
// max 64kb for each batch
func syncBatch(conn *net.TCPConn, startSeq int64) {
	l.Info().Msgf("persistence::syncBach start")

	batchSyncing = true
	defer func() {
		batchSyncing = false
	}()

	//
	fpath := getActiveBinlogFile()
	exist, err := util.PathExists(fpath)
	if err != nil {
		l.Err(err).Msgf("persistence::catchup %s", err.Error())
		return
	}
	if !exist {
		return
	}
	f, err := os.OpenFile(fpath, os.O_RDONLY, 0664)
	if err != nil {
		l.Err(err).Msgf("persistence::catchup %s", err.Error())
		return
	}
	defer func() {
		if f != nil {
			f.Close()
		}
	}()

	// locate via seq
	// TODO: to use index
	var seq int64
	var binlogs []byte
	for {
		buf, err := parseOne(f, &seq)
		if err != nil {
			if err != io.EOF {
				l.Panic().Msg(err.Error())
			}
			break
		}
		if seq < startSeq {
			continue
		}
		binlogs = append(binlogs, buf...)
		if len(binlogs) > batchSyncSize { // limit batch size
			break
		}
	}

	flrSeq, err := doSync(conn, binlogs)
	if err != nil {
		l.Err(err).Msgf("persistence::syncBach %s", err.Error())
		l.Warn().Msgf("persistence::syncBach end with error")
	}
	l.Info().Msgf("persistence::syncBach end")

	if flrSeq < binlogSeq { // continue
		syncBatch(conn, flrSeq+1)
	}
}

// Sync to follower(s)
func doSync(conn *net.TCPConn, binlogs []byte) (flrSeq int64, err error) {
	// icmd
	bbuf := new(bytes.Buffer)
	binary.Write(bbuf, binary.LittleEndian, consts.RPC_SYNC_BINLOG)

	// binlog
	buf := bbuf.Bytes()
	buf = append(buf, binlogs...)

	// encode
	msg, err := proto.EncodeB(buf)
	if err != nil {
		l.Err(err).Msgf("persistence::sync encode: %s", err.Error())
		return flrSeq, err
	}

	// send
	_, err = conn.Write(msg)
	if err != nil {
		l.Err(err).Msgf("persistence::sync send: %s", err.Error())
		return flrSeq, err
	}

	// get follower's seq
	reader := bufio.NewReader(conn)
	res, err := proto.Decode(reader)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msgf("persistence::sync res: %s", err.Error())
		}
		return flrSeq, err
	}

	if len(res) < resp.RESPONSE_HEADER_SIZE {
		return flrSeq, err
	}

	bbuf.Reset()

	bbuf = bytes.NewBuffer(res[resp.RESPONSE_HEADER_SIZE:])
	err = binary.Read(bbuf, binary.LittleEndian, &flrSeq)
	if err != nil {
		l.Err(err).Msgf("persistence::sync read res: %s", err.Error())
		return flrSeq, err
	}

	return flrSeq, nil
}

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
