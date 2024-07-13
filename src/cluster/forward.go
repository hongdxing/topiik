/***
** author: duan hongxing
** data: 4 Jul 2024
** desc:
**
**/

package cluster

import (
	"bufio"
	"io"
	"net"
	"topiik/internal/proto"
	"topiik/internal/util"
)

// cache Tcp Conn from Controller to Workers
var tcpMap = make(map[string]*net.TCPConn)

func Forward(msg []byte) []byte {
	if len(clusterInfo.Workers) == 0 {
		return []byte{}
	}
	var err error
	// TODO: find worker base on key partition, and get LeaderWorkerId
	// and then get Address of Worker

	var targetWorker NodeSlim
	for _, worker := range clusterInfo.Workers {
		targetWorker = worker
		break
	}

	conn, ok := tcpMap[targetWorker.Id]
	if !ok {
		conn, err = util.PreapareSocketClient(targetWorker.Address)
		if err != nil {
			return []byte{} // TODO: should retry
		}
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
