/***
** author: duan hongxing
** date: 26 Jun 2024
** desc:
**
**/

package persistence

import (
	"fmt"
	"sync"
	"time"
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/memo"
)

var persistenceTicker *time.Ticker
var persistenceWG sync.WaitGroup
var quit chan struct{}

func Persist(serverConfig config.ServerConfig) {
	persistenceTicker = time.NewTicker(time.Duration(serverConfig.SaveMillis) * time.Millisecond)
	quit = make(chan struct{})

	// Start compress routine
	go compress()

	for {
		select {
		case <-persistenceTicker.C:
			persistenceWG.Add(1)
			doPersist()
		case <-quit:
			persistenceTicker.Stop()
			return
		}
	}
}

func doPersist() {
	defer persistenceWG.Done()

	if memo.MemMap[consts.PERSISTENT_BUF_QUEUE] == nil {
		return
	}
	for ele := memo.MemMap[consts.PERSISTENT_BUF_QUEUE].Lst.Back(); ele != nil; ele = ele.Prev() {
		//result = append(result, ele.Value.(string))
		//eleToBeRemoved = append(eleToBeRemoved, ele)
		fmt.Printf("%b", ele.Value)
		memo.MemMap[consts.PERSISTENT_BUF_QUEUE].Lst.Remove(ele)
	}
}
