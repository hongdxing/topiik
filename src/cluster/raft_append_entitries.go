/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package cluster

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
)

var ticker *time.Ticker
var quit chan struct{}
var wgAppend sync.WaitGroup

// indicate metadata changed on controller Leader, need to sync to Follower(s)
var controllerPendingAppend = make(map[string]string) // the node id of controller wating for append
var workerPendingAppend = make(map[string]string)
var partitionPendingAppend = make(map[string]string)

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
			for _, controller := range controllerMap {
				if controller.Id == nodeInfo.Id {
					continue
				}
				wgAppend.Add(1)
				go send(controller.Address2, controller.Id, &dialErrorCounter)
			}
			wgAppend.Wait()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func send(address string, controllerId string, dialErrorCounter *int) string {
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

	line := RPC_APPENDENTRY + consts.SPACE
	if _, ok := controllerPendingAppend[controllerId]; ok {
		line += "CONTROLLER "
		buf, _ := json.Marshal(controllerMap)
		line += string(buf)
	} else if _, ok := workerPendingAppend[controllerId]; ok {
		line += "WORKER "
	} else if _, ok := partitionPendingAppend[controllerId]; ok {
		line += "PARTITION "
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
		return ""
	}

	/*
		reader := bufio.NewReader(conn)
		buf, err := proto.Decode(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("rpc_append_entries::send %s\n", err)
			}
		}
		fmt.Printf("rpc_append_entries::send %s\n", buf)
		return string(buf[4:])*/

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("rpc_append_entries::send %s\n", err)
	} else {
		delete(controllerPendingAppend, controllerId)
		//fmt.Println(string(buf[:n]))
	}
	return string(buf[:n])

}
