// author: Duan Hongxing
// date: 18 Sep, 2024

package persistence

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/rorre"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

// cache tcp conn to persistor
var pstConn *net.TCPConn

// conn cache from leader to follower
var followerConnCache = make(map[string]*net.TCPConn)

var msgQ list.List
var lock sync.Mutex
var dequeueTicker time.Ticker

// enqueue msg pending persist to persistor server
func Enqueue(msg []byte) {
	lock.Lock()
	defer lock.Unlock()
	msgQ.PushFront(msg)
}

// dequeue msg sync to persistor server
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
	_, err := syncToPersistor(binlogs)
	if err != nil {
		switch err.(type) {
		case *rorre.SoketError:
			if pstConn != nil {
				pstConn.Close()
			}
			pstConn = nil
		}
	}
}

// Sync to follower(s)
func syncToPersistor(binlogs []byte) (flrSeq int64, err error) {
	// icmd
	bbuf := new(bytes.Buffer)
	binary.Write(bbuf, binary.LittleEndian, consts.RPC_SYNC_BINLOG)

	// append partition id
	ptnId := node.GetNodeInfo().PntId
	if len(ptnId) <= 0 {
		l.Warn().Msgf("persistence::syncToPersistor Invalid partition Id")
		return
	}

	// binlog
	buf := bbuf.Bytes()                 //rpc type
	buf = append(buf, []byte(ptnId)...) //ptn id
	buf = append(buf, binlogs...)       //binlog

	// encode
	msg, err := proto.EncodeB(buf)
	if err != nil {
		l.Err(err).Msgf("persistence::sync encode: %s", err.Error())
		return flrSeq, err
	}

	// send
	_, err = pstConn.Write(msg)
	if err != nil {
		l.Err(err).Msgf("persistence::sync send: %s", err.Error())
		return flrSeq, &rorre.SoketError{}
	}

	// get follower's seq
	reader := bufio.NewReader(pstConn)
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

// sync to partition follower(s)
func SyncFollower(flrs []node.NodeSlim, binlogs []byte) (err error) {
	// for each follower
	for _, flr := range flrs {
		_, err := doSyncFollower(flr, binlogs)
		if err != nil {
			switch err.(type) {
			case *rorre.SoketError: // if is SocketError, close the conn and delete cache
				if conn, ok := followerConnCache[flr.Id]; ok {
					conn.Close()
					delete(followerConnCache, flr.Id)
				}
			}
		}
	}
	return nil
}

func doSyncFollower(nd node.NodeSlim, binlogs []byte) (flrSeq int64, err error) {
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

	var conn *net.TCPConn
	if cached, ok := followerConnCache[nd.Id]; ok {
		conn = cached
	} else {
		conn, err := util.PreapareSocketClient(nd.Addr2)
		if err != nil {
			l.Err(err).Msgf("persistence::doSyncFollower conn: %s", err.Error())
		}
		followerConnCache[nd.Id] = conn
	}

	// send
	_, err = conn.Write(msg)
	if err != nil {
		l.Err(err).Msgf("persistence::sync send: %s", err.Error())
		return flrSeq, &rorre.SoketError{}
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
