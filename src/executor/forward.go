/***
** author: duan hongxing
** data: 4 Jul 2024
** desc:
**
**/

package executor

import (
	"bufio"
	"errors"
	"hash/crc32"
	"io"
	"net"
	"topiik/cluster"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

/* Conn cache from Controller to Workers */
var ctlwkrConnCache = make(map[string]*net.TCPConn)

func forwardByKey(key string, msg []byte) []byte {
	if len(cluster.GetClusterInfo().Wkrs) == 0 {
		return resp.ErrorResponse(errors.New(resp.RES_NO_ENOUGH_WORKER))
	}
	if len(cluster.GetPartitionInfo().PtnMap) == 0 {
		return resp.ErrorResponse(errors.New(resp.RES_NO_PARTITION))
	}
	var err error
	// find worker base on key partition, and get LeaderWorkerId
	targetWorker := getWorker(key)
	if len(targetWorker.Id) == 0 {
		l.Warn().Msg("forward::Forward no slot available")
		return resp.ErrorResponse(errors.New(resp.RES_NO_PARTITION))
	}

	conn, ok := ctlwkrConnCache[targetWorker.Id]
	if !ok {
		//conn, err = util.PreapareSocketClientWithPort(targetWorker.Addr, CONTROLLER_FORWORD_PORT)
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return resp.ErrorResponse(errors.New(resp.RES_CONN_RESET))
		}
		ctlwkrConnCache[targetWorker.Id] = conn
	}
	// Send
	_, err = conn.Write(msg)
	if err != nil {
		l.Err(err).Msg(err.Error())
		if _, ok = ctlwkrConnCache[targetWorker.Id]; ok {
			l.Warn().Msgf("forward::Forward remove tcp cache of worker %s", targetWorker.Id)
			delete(ctlwkrConnCache, targetWorker.Id)
		}
		// try reconnect
		targetWorker := getWorker(key) // the worker may changed because of Worker Leader fail
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			l.Err(err).Msg(err.Error())
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

func getWorker(key string) (worker node.NodeSlim) {
	var keyHash = crc32.Checksum([]byte(key), crc32.IEEETable)
	keyHash = keyHash % consts.SLOTS
	//fmt.Printf("key hash %v\n", keyHash)
	for _, partition := range cluster.GetPartitionInfo().PtnMap {
		for _, slot := range partition.Slots {
			if slot.From <= uint16(keyHash) && slot.To >= uint16(keyHash) {
				worker = cluster.GetClusterInfo().Wkrs[partition.LeaderNodeId]
				break
			}
		}
		if len(worker.Id) > 0 {
			break
		}
	}
	return worker
}

func forwardByWorker(targetWorker node.NodeSlim, msg []byte) []byte {
	var err error
	conn, ok := ctlwkrConnCache[targetWorker.Id]
	if !ok {
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return resp.ErrorResponse(errors.New(resp.RES_CONN_RESET))
		}
		ctlwkrConnCache[targetWorker.Id] = conn
	}
	// Send
	_, err = conn.Write(msg)
	if err != nil {
		l.Err(err).Msg(err.Error())
		if _, ok = ctlwkrConnCache[targetWorker.Id]; ok {
			l.Warn().Msgf("forward::Forward remove tcp cache of worker %s", targetWorker.Id)
			delete(ctlwkrConnCache, targetWorker.Id)
		}
		// try reconnect
		//targetWorker := getWorker(key) // the worker may changed because of Worker Leader fail
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			l.Err(err).Msg(err.Error())
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
