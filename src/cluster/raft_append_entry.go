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
	"topiik/internal/proto"
	"topiik/internal/util"
)

var ticker *time.Ticker
var cluUpdCh chan struct{}
var ptnUpdCh chan struct{}

//var wgAppend sync.WaitGroup

// indicate metadata changed on controller Leader, need to sync to Follower(s)
//var clusterMetadataPendingAppend = make(map[string]string)   // the controller id, id
//var partitionMetadataPendingAppend = make(map[string]string) // the controller id, id

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
			//go appendWorkers(&dialErrorCounter)
			//appendControllers(&dialErrorCounter)
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
			if controller.Id == nodeInfo.Id {
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
			if controller.Id == nodeInfo.Id {
				continue
			}
			send(controller.Addr2, controller.Id, buf)
		}
	}
}

func appendHeartbeat() {
	for _, controller := range clusterInfo.Ctls {
		if controller.Id == nodeInfo.Id {
			continue
		}
		send(controller.Addr2, controller.Id, []byte{})
	}
	for _, worker := range clusterInfo.Wkrs {
		if worker.Id == nodeInfo.Id {
			continue
		}
		send(worker.Addr2, worker.Id, []byte{})
	}
}

/*
func appendControllers(dialErrorCounter *int) {
	for _, controller := range clusterInfo.Ctls {
		if controller.Id == nodeInfo.Id {
			continue
		}
		send(true, controller.Addr2, controller.Id, dialErrorCounter)
	}
}

func appendWorkers(dialErrorCounter *int) {
	for _, worker := range clusterInfo.Wkrs {
		if worker.Id == nodeInfo.Id {
			continue
		}
		//wgAppend.Add(1)
		send(false, worker.Addr2, worker.Id, dialErrorCounter)
		//wgAppend.Wait()
	}
}
*/

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
	binary.Write(byteBuf, binary.LittleEndian, RPC_APPENDENTRY)
	rpcBuf = append(rpcBuf, byteBuf.Bytes()...)

	if len(data) > 0 {
		rpcBuf = append(rpcBuf, data...)
	}

	if len(rpcBuf) == 1 { // means no data, then append controller's addr
		byteBuf.Reset()
		binary.Write(byteBuf, binary.LittleEndian, ENTRY_TYPE_DEFAULT)
		rpcBuf = append(rpcBuf, byteBuf.Bytes()...)
		rpcBuf = append(rpcBuf, []byte(nodeInfo.Addr)...)
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

/*
func send1(isController bool, destAddr string, nodeId string, dialErrorCounter *int) string {
	defer func() {
		*dialErrorCounter++
		if *dialErrorCounter >= 10000 {
			*dialErrorCounter = 0
		}
		//wgAppend.Done()
	}()

	var err error
	var conn *net.TCPConn

	if v, ok := connCache[nodeId]; ok {
		conn = v
	}
	if conn == nil {
		conn, err = util.PreapareSocketClient(destAddr)
		if err != nil {
			if *dialErrorCounter%50 == 0 {
				l.Err(err)
			}
			return ""
		}
		connCache[nodeId] = conn
	}

	var cmdBytes []byte
	var byteBuf = new(bytes.Buffer) // int to byte byte buf
	// 1 bytes of command + 1 byte of entry type + the entry
	binary.Write(byteBuf, binary.LittleEndian, RPC_APPENDENTRY)
	cmdBytes = append(cmdBytes, byteBuf.Bytes()...)

	if isController {
		if _, ok := clusterMetadataPendingAppend[nodeId]; ok { // if metadata pending
			byteBuf.Reset()
			buf, err := json.Marshal(clusterInfo)
			if err != nil {
				l.Err(err).Msg(err.Error())
			} else {
				binary.Write(byteBuf, binary.LittleEndian, ENTRY_TYPE_METADATA)
				cmdBytes = append(cmdBytes, byteBuf.Bytes()...)
				cmdBytes = append(cmdBytes, buf...)
			}
		} else if _, ok := partitionMetadataPendingAppend[nodeId]; ok { // if partition pending
			byteBuf.Reset()
			buf, err := json.Marshal(partitionInfo)
			if err != nil {
				l.Err(err).Msg(err.Error())
			} else {
				binary.Write(byteBuf, binary.LittleEndian, ENTRY_TYPE_PARTITION)
				cmdBytes = append(cmdBytes, byteBuf.Bytes()...)
				cmdBytes = append(cmdBytes, buf...)
			}
		}
	}

	if len(cmdBytes) == 1 { // means no data, then append controller's addr
		byteBuf.Reset()
		binary.Write(byteBuf, binary.LittleEndian, ENTRY_TYPE_DEFAULT)
		cmdBytes = append(cmdBytes, byteBuf.Bytes()...)
		cmdBytes = append(cmdBytes, []byte(nodeInfo.Addr)...)
	}

	// Enocde
	data, err := proto.EncodeB(cmdBytes)
	if err != nil {
		l.Err(err)
	}

	// Send
	_, err = conn.Write(data)
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
	// remove the pending conroller id from Pending map
	delete(clusterMetadataPendingAppend, nodeId)
	delete(partitionMetadataPendingAppend, nodeId)
	return string(buf)
}
*/
