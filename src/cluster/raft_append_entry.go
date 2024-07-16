/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package cluster

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
)

var ticker *time.Ticker
var quit chan struct{}

//var wgAppend sync.WaitGroup

// indicate metadata changed on controller Leader, need to sync to Follower(s)
var clusterMetadataPendingAppend = make(map[string]string) // the controller id, id

var connCache = make(map[string]*net.TCPConn)

/***
** Controller issues AppendEntries RPCs to replicate metadata to follower,
** or send heartbeats (AppendEntries RPCs that carry no data)
**
**/
func AppendEntries() {
	ticker = time.NewTicker(200 * time.Millisecond)
	quit = make(chan struct{})

	dialErrorCounter := 0
	for {
		select {
		case <-ticker.C:
			for _, controller := range clusterInfo.Controllers {
				if controller.Id == nodeInfo.Id {
					continue
				}
				//wgAppend.Add(1)
				send(controller.Address2, controller.Id, &dialErrorCounter)
				//wgAppend.Wait()
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func send(address string, controllerId string, dialErrorCounter *int) string {
	defer func() {
		*dialErrorCounter++
		//wgAppend.Done()
	}()

	var err error
	var conn *net.TCPConn

	if v, ok := connCache[controllerId]; ok {
		conn = v
	}
	if conn == nil {
		conn, err = util.PreapareSocketClient(address)
		if err != nil {
			if *dialErrorCounter%50 == 0 {
				fmt.Println(err)
			}
			return ""
		}
		connCache[controllerId] = conn
	}
	//defer conn.Close()

	line := RPC_APPENDENTRY + consts.SPACE
	if _, ok := clusterMetadataPendingAppend[controllerId]; ok {
		line += "METADATA "
		buf, _ := json.Marshal(clusterInfo)
		line += string(buf)
	}

	// Enocde
	data, err := proto.Encode(line)
	if err != nil {
		fmt.Println(err)
	}

	// Send
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println(err)
		if conn, ok := connCache[controllerId]; ok {
			conn.Close()
			conn = nil
			fmt.Println("raft_append_entries::send Delete connCache")
			delete(connCache, controllerId)
		}

		return ""
	}

	reader := bufio.NewReader(conn)
	buf, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			fmt.Printf("rpc_append_entries::send %s\n", err)
		}
	}
	// remove the pending conroller id from Pending map
	delete(clusterMetadataPendingAppend, controllerId)
	return string(buf[4:])
}
