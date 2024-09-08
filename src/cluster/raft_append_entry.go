//author: Duan Hongxing
//data: 3 Jul, 2024

package cluster

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
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

const ptnTickerDur int = 500
const ptnLeaderDownMaxMills int = 3000

var hbTicker *time.Ticker  // heartbeat ticker
var ptnTicker *time.Ticker // partition ticker
// var cluUpdCh chan struct{}

var connCache = make(map[string]*net.TCPConn)

// Controller issues AppendEntries RPCs to replicate metadata to follower,
// or send heartbeats (AppendEntries RPCs that carry no data)
func AppendEntries() {
	hbTicker = time.NewTicker(200 * time.Millisecond)
	ptnTicker = time.NewTicker(time.Duration(ptnTickerDur) * time.Millisecond)
	defer close(ptnUpdCh)
	defer close(ctlUpdCh)
	defer close(wrkUpdCh)

	//dialErrorCounter := 0 // this not 'thread' safe, but it's not important
	for {
		select {
		case <-hbTicker.C:
			appendHeartbeat()
		case <-ptnTicker.C:
			appendPtn()
		case <-ctlUpdCh:
			appendControllerInfo()
		case <-wrkUpdCh:
			appendWorkerInfo()
		case <-ptnUpdCh:
			appendPartitionInfo()
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

// sync controller info to controllers and workers
func appendControllerInfo() {
	var buf []byte
	var bbuf = new(bytes.Buffer)
	data, err := json.Marshal(controllerInfo)
	if err != nil {
		l.Err(err).Msg(err.Error())
	} else {
		binary.Write(bbuf, binary.LittleEndian, ENTRY_TYPE_CTL)
		buf = append(buf, bbuf.Bytes()...)
		buf = append(buf, data...)
		for _, controller := range controllerInfo.Nodes {
			if controller.Id == node.GetNodeInfo().Id {
				continue
			}
			send(controller.Addr2, controller.Id, buf)
		}
		for _, worker := range workerInfo.Nodes {
			if worker.Id == node.GetNodeInfo().Id {
				continue
			}
			send(worker.Addr2, worker.Id, buf)
		}
	}
}

// sync worker info to controllers
func appendWorkerInfo() {
	var buf []byte
	var bbuf = new(bytes.Buffer)
	data, err := json.Marshal(workerInfo)
	if err != nil {
		l.Err(err).Msg(err.Error())
	} else {
		binary.Write(bbuf, binary.LittleEndian, ENTRY_TYPE_WRK)
		buf = append(buf, bbuf.Bytes()...)
		buf = append(buf, data...)
		for _, controller := range controllerInfo.Nodes {
			if controller.Id == node.GetNodeInfo().Id {
				continue
			}
			send(controller.Addr2, controller.Id, buf)
		}
	}
}

func appendPartitionInfo() {
	var buf []byte
	var bbuf = new(bytes.Buffer)
	data, err := json.Marshal(partitionInfo)
	if err != nil {
		l.Err(err).Msg(err.Error())
	} else {
		binary.Write(bbuf, binary.LittleEndian, ENTRY_TYPE_PTNS)
		buf = append(buf, bbuf.Bytes()...)
		buf = append(buf, data...)
		for _, controller := range controllerInfo.Nodes {
			if controller.Id == node.GetNodeInfo().Id {
				continue
			}
			send(controller.Addr2, controller.Id, buf)
		}
	}
}

func appendHeartbeat() {
	for _, controller := range controllerInfo.Nodes {
		if controller.Id == node.GetNodeInfo().Id {
			continue
		}
		go send(controller.Addr2, controller.Id, []byte{})
	}

	for _, worker := range workerInfo.Nodes {
		if worker.Id == node.GetNodeInfo().Id {
			continue
		}
		go send(worker.Addr2, worker.Id, []byte{})
	}
}

// when partition leader not available, the accoumulated milli seconds retried
var ptnLeaderDownMills = make(map[string]int)

// healthcheck partition leader, and sync follower(s) to leader
func appendPtn() {
	for _, ptn := range partitionInfo.PtnMap {
		if ptn.LeaderNodeId == "" { // no leader yet
			electPtnLeader(ptn)
		} else {
			if ptnLeader, ok := workerInfo.Nodes[ptn.LeaderNodeId]; ok { // get the leader Worker
				for ndId, nd := range ptn.NodeSet {
					if wrk, ok := workerInfo.Nodes[ndId]; ok {
						nd.Id = ndId
						nd.Addr = wrk.Addr
						nd.Addr2 = wrk.Addr2
					}
				}

				var buf []byte
				bbuf := new(bytes.Buffer)
				binary.Write(bbuf, binary.LittleEndian, ENTRY_TYPE_PTN)
				buf = append(buf, bbuf.Bytes()...)
				//flrb, err := json.Marshal(followers)
				ptnb, err := json.Marshal(ptn)
				if err != nil {
					l.Err(err).Msgf(err.Error())
					continue
				}
				buf = append(buf, ptnb...)
				err = send(ptnLeader.Addr2, ptnLeader.Id, buf)
				if err != nil {
					//l.Err(err).Msg(err.Error())
					var ne net.Error
					if errors.As(err, &ne) { // if the leader not available from controller
						l.Warn().Msgf("worker %s not accessible", ptnLeader.Addr2)
						ptnLeaderDownMills[ptn.Id] += ptnTickerDur
						if ptnLeaderDownMills[ptn.Id] >= ptnLeaderDownMaxMills { // to elect new partition leader
							electPtnLeader(ptn)
						}
					}
				} else {
					ptnLeaderDownMills[ptn.Id] = 0
				}
			}
		}
	}
}

func electPtnLeader(ptn *node.Partition) {
	seqMap := make(map[string]int64)
	var wg sync.WaitGroup
	var newLeaderId string

	// set init value to first node id
	for _, nd := range ptn.NodeSet {
		newLeaderId = nd.Id
		break
	}

	if len(ptn.NodeSet) > 1 {
		// notice here not execlude the leader, in case it's recovered and give it last chance
		for ndId := range ptn.NodeSet {
			wg.Add(1)
			if wrk, ok := workerInfo.Nodes[ndId]; ok {
				go getWorkerBinlogSeq(wrk.Id, wrk.Addr2, &wg, &seqMap)
			} else {
				wg.Done()
			}
		}
		wg.Wait()

		// get worker that have max seq
		var maxSeq int64
		for ndId, seq := range seqMap {
			if seq > maxSeq {
				newLeaderId = ndId
				maxSeq = seq
			}
		}
	}

	// l.Info().Msgf("new LeaderId: %s", newLeaderId)
	if _, ok := workerInfo.Nodes[newLeaderId]; ok {
		l.Info().Msgf("Partition %s has new leader: %s", ptn.Id, newLeaderId)
		ptn.LeaderNodeId = newLeaderId
		//ptnUpdCh <- struct{}{} // sync partition to followers
	}
}

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
