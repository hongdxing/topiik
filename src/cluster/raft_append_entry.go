/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package cluster

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"time"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
)

var ticker *time.Ticker
var cluUpdCh chan struct{}
var ptnUpdCh chan struct{}

var connCache = make(map[string]*net.TCPConn)

/***
** Controller issues AppendEntries RPCs to replicate metadata to follower,
** or send heartbeats (AppendEntries RPCs that carry no data)
**
**/
func AppendEntries() {
	ticker = time.NewTicker(200 * time.Millisecond)
	cluUpdCh = make(chan struct{})
	ptnUpdCh = make(chan struct{})

	//dialErrorCounter := 0 // this not 'thread' safe, but it's not important
	for {
		select {
		case <-ticker.C:
			appendHeartbeat()
		case <-cluUpdCh:
			appendClusterInfo()
		case <-ptnUpdCh:
			appendPartitionInfo()
		}
	}
}

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
		for _, controller := range clusterInfo.Ctls {
			if controller.Id == node.GetNodeInfo().Id {
				continue
			}
			send(controller.Addr2, controller.Id, buf)
		}
	}
}

func appendPartitionInfo() {
	var buf []byte
	var byteBuf = new(bytes.Buffer)
	data, err := json.Marshal(partitionInfo)
	if err != nil {
		l.Err(err).Msg(err.Error())
	} else {
		binary.Write(byteBuf, binary.LittleEndian, ENTRY_TYPE_PARTITION)
		buf = append(buf, byteBuf.Bytes()...)
		buf = append(buf, data...)
		for _, controller := range clusterInfo.Ctls {
			if controller.Id == node.GetNodeInfo().Id {
				continue
			}
			send(controller.Addr2, controller.Id, buf)
		}
	}
}

func appendHeartbeat() {
	for _, controller := range clusterInfo.Ctls {
		if controller.Id == node.GetNodeInfo().Id {
			continue
		}
		send(controller.Addr2, controller.Id, []byte{})
	}
	for _, worker := range clusterInfo.Wkrs {
		if worker.Id == node.GetNodeInfo().Id {
			continue
		}
		send(worker.Addr2, worker.Id, []byte{})
	}
}

func send(destAddr string, nodeId string, data []byte) string {
	var err error
	var conn *net.TCPConn

	if v, ok := connCache[nodeId]; ok {
		conn = v
	}
	if conn == nil {
		conn, err = util.PreapareSocketClient(destAddr)
		if err != nil {
			return ""
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

	if len(rpcBuf) == 1 { // means no data, then append controller's addr
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
			l.Warn().Msg("raft_append_entries::send Delete connCache")
			delete(connCache, nodeId)
		}

		return ""
	}

	reader := bufio.NewReader(conn)
	buf, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			l.Err(err).Msgf("rpc_append_entries::send %s\n", err)
		}
	}
	return string(buf)
}
