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
)

var hbTicker *time.Ticker // heartbeat ticker

var connCache = make(map[string]*net.TCPConn)

// Controller issues AppendEntries RPCs to replicate metadata to follower,
// or send heartbeats (AppendEntries RPCs that carry no data)
func AppendEntries(ptn Partition) {
	hbTicker = time.NewTicker(200 * time.Millisecond)
	defer close(ptnUpdCh)

	//dialErrorCounter := 0 // this not 'thread' safe, but it's not important
	for {
		select {
		case <-hbTicker.C:
			appendHeartbeat(ptn)
		case <-ptnUpdCh:
			appendPtnInfo()
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
func appendPtnInfo() {
	var buf []byte
	var bbuf = new(bytes.Buffer)
	data, err := json.Marshal(partitionInfo)
	if err != nil {
		l.Err(err).Msg(err.Error())
	} else {
		binary.Write(bbuf, binary.LittleEndian, ENTRY_TYPE_WRKGRP)
		buf = append(buf, bbuf.Bytes()...)
		buf = append(buf, data...)
		for _, ptn := range partitionInfo.Ptns {
			for _, worker := range ptn.Nodes {
				if worker.Id == node.GetNodeInfo().Id {
					continue
				}
				send(worker.Addr2, worker.Id, buf)
			}
		}
	}
}

func appendHeartbeat(ptn Partition) {
	for _, controller := range ptn.Nodes {
		if controller.Id == node.GetNodeInfo().Id {
			continue
		}
		go send(controller.Addr2, controller.Id, []byte{})
	}

	for _, worker := range ptn.Nodes {
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
