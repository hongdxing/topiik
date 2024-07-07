/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package cluster

import (
	"fmt"
	"sync"
	"time"
	"topiik/internal/proto"
	"topiik/internal/util"
)

var ticker *time.Ticker
var quit chan struct{}
var wgAppend sync.WaitGroup

/***
** Controller issues AppendEntries RPCs to replicate metadata to deputy,
** or send heartbeats (AppendEntries RPCs that carry no data)
**
**	Parameters:
**	- addresses: actually only one address of Chief Officer
**/
func AppendEntries() {
	ticker = time.NewTicker(200 * time.Millisecond)
	quit = make(chan struct{})

	var controllerAddrs []string
	dialErrorCounter := 0
	for {
		select {
		case <-ticker.C:
			clear(controllerAddrs)
			for _, v := range controllerMap {
				controllerAddrs = append(controllerAddrs, v.Address2)
			}
			for _, address := range controllerAddrs {
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

	conn, err := util.PreapareSocketClient(address)
	if err != nil {
		if *dialErrorCounter%50 == 0 {
			fmt.Println(err)
		}
		return ""
	}
	defer conn.Close()

	line := RPC_APPENDENTRY + " "

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
