/***
**
**
**
**
**
**/

package util

import (
	"fmt"
	"net"
)

func PreapareSocketClient(addr string) (tcpConn *net.TCPConn, err error) {
	rAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
	}
	tcpConn, err = net.DialTCP("tcp", nil, rAddr)
	return tcpConn, err
}

func PreapareSocketClientWithPort(addr string, clientPort string) (*net.TCPConn, error) {
	rAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
	}
	lAddr, err := net.ResolveTCPAddr("tcp", ":"+clientPort)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
	}
	tcpConn, err := net.DialTCP("tcp", lAddr, rAddr)
	if err != nil {
		fmt.Printf("PreapareSocketClientWithPort() %s\n", err.Error())
	}
	return tcpConn, err
}
