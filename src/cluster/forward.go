/***
** author: duan hongxing
** data: 4 Jul 2024
** desc:
**
**/

package cluster

import (
	"bufio"
	"errors"
	"hash/crc32"
	"io"
	"net"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/resp"
)

// cache Tcp Conn from Controller to Workers
var tcpMap = make(map[string]*net.TCPConn)

func Forward(key string, msg []byte) []byte {
	if len(clusterInfo.Wkrs) == 0 {
		//res, _ := proto.Encode("")
		//return res
		return resp.ErrorResponse(errors.New(resp.RES_NO_ENOUGH_WORKER))
	}
	if len(partitionInfo) == 0 {
		return resp.ErrorResponse(errors.New(resp.RES_NO_PARTITION))
	}
	var err error
	// find worker base on key partition, and get LeaderWorkerId
	var targetWorker Worker
	var keyHash = crc32.Checksum([]byte(key), crc32.IEEETable)
	keyHash = keyHash % SLOTS
	//fmt.Printf("key hash %v\n", keyHash)

	for _, partition := range partitionInfo {
		for _, slot := range partition.Slots {
			if slot.From <= uint16(keyHash) && slot.To >= uint16(keyHash) {
				targetWorker = clusterInfo.Wkrs[partition.LeaderNodeId]
				break
			}
		}
		if len(targetWorker.Id) > 0 {
			break
		}
	}
	if len(targetWorker.Id) == 0 {
		tLog.Warn().Msg("forward::Forward no slot available")
		return resp.ErrorResponse(errors.New(resp.RES_NO_PARTITION))
	}

	conn, ok := tcpMap[targetWorker.Id]
	if !ok {
		//conn, err = util.PreapareSocketClientWithPort(targetWorker.Addr, CONTROLLER_FORWORD_PORT)
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			tLog.Err(err).Msg(err.Error())
			return resp.ErrorResponse(errors.New(resp.RES_CONN_RESET))
		}
		tcpMap[targetWorker.Id] = conn
	}
	// Send
	_, err = conn.Write(msg)
	if err != nil {
		tLog.Err(err).Msg(err.Error())
		if _, ok = tcpMap[targetWorker.Id]; ok {
			tLog.Warn().Msgf("Forward() delete worker %s from tcpMap", targetWorker.Id)
			delete(tcpMap, targetWorker.Id)
		}
		// try reconnect
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			tLog.Err(err).Msg(err.Error())
			return resp.ErrorResponse(errors.New(resp.RES_CONN_RESET))
		}
		conn.Write(msg)
	}

	reader := bufio.NewReader(conn)
	responseBytes, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			return resp.ErrorResponse(errors.New(resp.RES_CONN_RESET))
		}
	}
	return responseBytes
}
