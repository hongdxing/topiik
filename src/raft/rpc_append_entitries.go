/***
* author: duan hongxing
* date: 21 Jun 2024
* desc:
 */

package raft

import (
	"fmt"
	"net"
	"sync"
	"time"
	"topiik/internal/proto"
)

var ticker *time.Ticker
var quit chan struct{}
var wgAppend sync.WaitGroup

/***
* leader issues AppendEntries RPCs to replicate log entries to followers,
* or send heartbeats (AppendEntries RPCs that carry no log entries)
 */
func AppendEntries(addresses []string) {
	ticker = time.NewTicker(200 * time.Millisecond)
	quit = make(chan struct{})
	go doAppendEntries(addresses)
}

func doAppendEntries(addresses []string) {
	dialErrorCounter := 0
	for {
		select {
		case <-ticker.C:
			for _, address := range addresses {
				wgAppend.Add(1)
				go send(address, &dialErrorCounter)
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func send(address string, dialErrorCounter *int) string {
	defer func() {
		*dialErrorCounter++
		wgAppend.Done()
	}()
	tcpServer, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		if *dialErrorCounter%50 == 0 {
			println("ResolveTCPAddr failed:", err.Error())
		}

	}
	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		if *dialErrorCounter%50 == 0 {
			fmt.Println(err)
		}
		return ""
	}
	defer conn.Close()

	//line := command.APPEND_ENTRY + " "
	line := "APPEND_ENTRY" + " "

	// Enocde
	data, err := proto.Encode(line)
	if err != nil {
		fmt.Println(err)
	}

	// Send
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("rpc_append_entries::send %s\n", err)
	} else {
		//fmt.Println(string(buf[:n]))
	}
	return string(buf[:n])
}
