//author: Duan HongXing
//data: 4 Jul, 2024

package shared

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

// Conn cache from leader to leader
var ctlwkrConnCache = make(map[string]*net.TCPConn)

func ForwardByKey(key []byte, msg []byte) []byte {
	if len(cluster.GetPartitionInfo().PtnMap) == 0 {
		return resp.ErrResponse(errors.New(resp.RES_NO_PARTITION))
	}
	var err error
	// find worker base on key partition, and get LeaderWorkerId
	targetWorker := getLeaderNode(key)
	if len(targetWorker.Id) == 0 {
		l.Warn().Msg("forward::Forward no slot available")
		return resp.ErrResponse(errors.New(resp.RES_NO_WORKER))
	}

	conn, ok := ctlwkrConnCache[targetWorker.Id]
	if !ok {
		//conn, err = util.PreapareSocketClientWithPort(targetWorker.Addr, CONTROLLER_FORWORD_PORT)
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return resp.ErrResponse(errors.New(resp.RES_CONN_RESET))
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
		targetWorker := getLeaderNode(key) // the worker may changed because of Worker Leader fail
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return resp.ErrResponse(errors.New(resp.RES_CONN_RESET))
		}
		conn.Write(msg)
	}

	reader := bufio.NewReader(conn)
	responseBytes, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			return resp.ErrResponse(errors.New(resp.RES_CONN_RESET))
		}
	}
	return responseBytes
}

func getLeaderNode(key []byte) (node node.NodeSlim) {
	var keyHash = crc32.Checksum(key, crc32.IEEETable)
	keyHash = keyHash % consts.SLOTS
	node = cluster.GetNodeByKeyHash(uint16(keyHash))
	return node
}

func ForwardByWorker(targetWorker node.NodeSlim, msg []byte) []byte {
	var err error
	conn, ok := ctlwkrConnCache[targetWorker.Id]
	if !ok {
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return resp.ErrResponse(errors.New(resp.RES_CONN_RESET))
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
			return resp.ErrResponse(errors.New(resp.RES_CONN_RESET))
		}
		conn.Write(msg)
	}

	reader := bufio.NewReader(conn)
	responseBytes, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			return resp.ErrResponse(errors.New(resp.RES_CONN_RESET))
		}
	}
	return responseBytes
}
