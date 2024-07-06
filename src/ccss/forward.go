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

// cache Tcp Conn from Capital to Sailors
var tcpMap = make(map[string]*net.TCPConn)

func Forward(msg []byte) []byte {
	if len(sailorMap) == 0 {
		return []byte{}
	}
	var err error
	// TODO: find sailor base on key partition, and get LeaderSailorId
	// and then get Address of Sailor

	var targetSailor Sailor
	for _, sailor := range sailorMap {
		targetSailor = sailor
		break
	}

	conn, ok := tcpMap[targetSailor.Id]
	if !ok {
		conn, err = util.PreapareSocketClient(targetSailor.Address)
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
