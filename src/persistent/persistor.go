/***
** author: duan hongxing
** date: 26 Jun 2024
** desc:
**
**/

package persistent

import (
	"fmt"
	"sync"
	"time"
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/shared"
)

var persistentTicker *time.Ticker
var persistentWG sync.WaitGroup
var quit chan struct{}

func Persist(serverConfig config.ServerConfig) {
	persistentTicker = time.NewTicker(time.Duration(serverConfig.SaveMillis) * time.Millisecond)
	quit = make(chan struct{})

	for {
		select {
		case <-persistentTicker.C:
			persistentWG.Add(1)
			doPersist()
		case <-quit:
			persistentTicker.Stop()
			return
		}
	}
}

func doPersist() {
	defer persistentWG.Done()

	if shared.MemMap[consts.PERSISTENT_BUF_QUEUE] == nil {
		return
	}
	for ele := shared.MemMap[consts.PERSISTENT_BUF_QUEUE].TList.Back(); ele != nil; ele = ele.Prev() {
		//result = append(result, ele.Value.(string))
		//eleToBeRemoved = append(eleToBeRemoved, ele)
		fmt.Printf("%b", []byte(ele.Value.(string)))
		shared.MemMap[consts.PERSISTENT_BUF_QUEUE].TList.Remove(ele)
	}
}
