/***
**
**
**
**
**
**/

package util

import (
	"net"
)

func PreapareSocketClient(addr string) (tcpConn *net.TCPConn, err error) {
	tcpServer, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
	}
	tcpConn, err = net.DialTCP("tcp", nil, tcpServer)
	return tcpConn, err
}
