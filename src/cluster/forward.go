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
		conn, err = util.PreapareSocketClientWithPort(targetWorker.Addr, CONTROLLER_FORWORD_PORT)
		//conn, err = util.PreapareSocketClient(targetWorker.Address)
		if err != nil {
			return []byte{} // TODO: should retry
		}
		tcpMap[targetWorker.Id] = conn
	}
	// Send
	_, err = conn.Write(msg)
	if err != nil {
		return []byte{} // TODO: should retry
	}

	reader := bufio.NewReader(conn)
	responseBytes, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			return []byte{}
		}
	}
	return responseBytes
}
