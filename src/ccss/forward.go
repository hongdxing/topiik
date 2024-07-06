/***
** author: duan hongxing
** data: 4 Jul 2024
** desc:
**
**/

package ccss

import (
	"bufio"
	"io"
	"net"
	"topiik/internal/proto"
	"topiik/internal/util"
)

// cache Tcp Conn from Capital to Salors
var tcpMap = make(map[string]*net.TCPConn)

func Forward(msg []byte) []byte {
	if len(salorMap) == 0 {
		return []byte{}
	}
	var err error
	// TODO: find salor base on key partition, and get LeaderSalorId
	// and then get Address of Salor

	var targetSalor Salor
	for _, salor := range salorMap {
		targetSalor = salor
		break
	}

	conn, ok := tcpMap[targetSalor.Id]
	if !ok {
		conn, err = util.PreapareSocketClient(targetSalor.Address)
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
