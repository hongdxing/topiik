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
	"io"
	"net"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/resp"
)

// cache Tcp Conn from Controller to Workers
var tcpMap = make(map[string]*net.TCPConn)

func Forward(msg []byte) []byte {
	if len(clusterInfo.Wkrs) == 0 {
		//res, _ := proto.Encode("")
		//return res
		return resp.ErrorResponse(errors.New(resp.RES_NO_ENOUGH_WORKER))
	}
	if len(partitionInfo) == 0 {
		return resp.ErrorResponse(errors.New(resp.RES_NO_PARTITION))
	}
	var err error
	// TODO: find worker base on key partition, and get LeaderWorkerId

	// and then get Address of Worker

	var targetWorker Worker
	for _, worker := range clusterInfo.Wkrs {
		targetWorker = worker
		break
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
