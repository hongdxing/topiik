//author: Duan Hongxing
//data: 3 Jul, 2024

package cluster

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"sync"
	"time"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

var hbTicker *time.Ticker // heartbeat ticker

var connCache = make(map[string]*net.TCPConn)

// Controller issues AppendEntries RPCs to replicate metadata to follower,
// or send heartbeats (AppendEntries RPCs that carry no data)
func AppendEntries(wrkGrp WorkerGroup) {
	hbTicker = time.NewTicker(200 * time.Millisecond)
	defer close(ptnUpdCh)
	defer close(wrkGrpUpdCh)

	//dialErrorCounter := 0 // this not 'thread' safe, but it's not important
	for {
		select {
		case <-hbTicker.C:
			appendHeartbeat(wrkGrp)
		case <-wrkGrpUpdCh:
			appendWrkGrpInfo()
		}
	}
}

/*
func appendClusterInfo() {
	var buf []byte
	var byteBuf = new(bytes.Buffer)
	data, err := json.Marshal(clusterInfo)
	if err != nil {
		l.Err(err).Msg(err.Error())
	} else {
		binary.Write(byteBuf, binary.LittleEndian, ENTRY_TYPE_METADATA)
		buf = append(buf, byteBuf.Bytes()...)
		buf = append(buf, data...)
		for _, controller := range controllerInfo.Nodes {
			if controller.Id == node.GetNodeInfo().Id {
				continue
			}
			send(controller.Addr2, controller.Id, buf)
		}
	}
}*/

// sync worker group info to workers
func appendWrkGrpInfo() {
	var buf []byte
	var bbuf = new(bytes.Buffer)
	data, err := json.Marshal(workerGroupInfo)
	if err != nil {
		l.Err(err).Msg(err.Error())
	} else {
		binary.Write(bbuf, binary.LittleEndian, ENTRY_TYPE_WRKGRP)
		buf = append(buf, bbuf.Bytes()...)
		buf = append(buf, data...)
		for _, wrkGrp := range workerGroupInfo.Groups {
			for _, worker := range wrkGrp.Nodes {
				if worker.Id == node.GetNodeInfo().Id {
					continue
				}
				send(worker.Addr2, worker.Id, buf)
			}
		}
	}
}

func appendHeartbeat(wrkGrp WorkerGroup) {
	for _, controller := range wrkGrp.Nodes {
		if controller.Id == node.GetNodeInfo().Id {
			continue
		}
		go send(controller.Addr2, controller.Id, []byte{})
	}

	for _, worker := range wrkGrp.Nodes {
		if worker.Id == node.GetNodeInfo().Id {
			continue
		}
		go send(worker.Addr2, worker.Id, []byte{})
	}
}

// when partition leader not available, the accoumulated milli seconds retried
var ptnLeaderDownMills = make(map[string]int)

var mu sync.Mutex

func send(destAddr string, nodeId string, data []byte) (err error) {
	var conn *net.TCPConn
	mu.Lock()
	defer mu.Unlock()
	if v, ok := connCache[nodeId]; ok {
		conn = v
	}
	if conn == nil {
		conn, err = util.PreapareSocketClient(destAddr)
		if err != nil {
			//l.Err(err).Msgf("cluster::send %s", err.Error())
			return err
		}
		connCache[nodeId] = conn
	}

	var rpcBuf []byte
	var byteBuf = new(bytes.Buffer) // int to byte byte buf
	// 1 bytes of command + 1 byte of entry type + the entry
	binary.Write(byteBuf, binary.LittleEndian, consts.RPC_APPENDENTRY)
	rpcBuf = append(rpcBuf, byteBuf.Bytes()...)

	if len(data) > 0 {
		rpcBuf = append(rpcBuf, data...)
	}

	/* means no data, then append controller's addr */
	if len(rpcBuf) == 1 {
		byteBuf.Reset()
		binary.Write(byteBuf, binary.LittleEndian, ENTRY_TYPE_DEFAULT)
		rpcBuf = append(rpcBuf, byteBuf.Bytes()...)
		rpcBuf = append(rpcBuf, []byte(node.GetNodeInfo().Addr)...)
	}

	// Enocde
	req, err := proto.EncodeB(rpcBuf)
	if err != nil {
		l.Err(err)
	}

	// Send
	_, err = conn.Write(req)
	if err != nil {
		l.Err(err)
		if conn, ok := connCache[nodeId]; ok {
			conn.Close()
			conn = nil
			l.Warn().Msg("rpc_append_entries::send Delete connCache")
			delete(connCache, nodeId)
		}

		return err
	}

	reader := bufio.NewReader(conn)
	_, err = proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			l.Err(err).Msgf("rpc_append_entries::send %s\n", err)
		}
	}
	return err
}

/*
* get binlog seq from worker
*
 */
func getWorkerBinlogSeq(ndId string, addr2 string, wg *sync.WaitGroup, seqMap *map[string]int64) {
	defer wg.Done()
	var conn *net.TCPConn
	var err error
	//var ok bool
	if c, ok := connCache[ndId]; ok {
		conn = c
	} else {
		conn, err = util.PreapareSocketClient(addr2)
		if err != nil {
			return
		}
		connCache[ndId] = conn
	}
	bbuf := new(bytes.Buffer)
	binary.Write(bbuf, binary.LittleEndian, consts.RPC_GET_BLSEQ)
	req := bbuf.Bytes()
	req, err = proto.EncodeB(req)
	if err != nil {
		return
	}
	_, err = conn.Write(req)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msgf("cluster::getWorkerBinlogSeq write %s", err.Error())
			return
		}
	}
	var seq int64
	reader := bufio.NewReader(conn)
	res, err := proto.Decode(reader)
	if err != nil {
		return
	}
	if len(res) > resp.RESPONSE_HEADER_SIZE {
		bbuf = bytes.NewBuffer(res[resp.RESPONSE_HEADER_SIZE:])
		err = binary.Read(bbuf, binary.LittleEndian, &seq)
		if err != nil {
			l.Err(err).Msgf("cluster::getWorkerBinlogSeq read %s", err.Error())
			return
		}
		l.Info().Msgf("worker %s seq is: %v", ndId, seq)
		(*seqMap)[ndId] = seq
	} else {
		l.Warn().Msgf("cluster::getWorkerBinlogSeq failed")
	}
}
