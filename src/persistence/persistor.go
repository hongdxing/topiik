/***
** author: duan hongxing
** date: 26 Jun 2024
** desc:
**
**/

package persistence

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"os"
	"topiik/executor"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/logger"
)

/*
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
*/

var l = logger.Get()
var newLineB = byte('\n')
var newLine = []byte{'\n'}

const maxCapacity int = 1024 * 1024 //

func Persist() {
	filePath := getCurrentLogFile()
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		l.Panic().Msg(err.Error())
	}
	defer file.Close()
	for {
		buf := <-executor.PersistenceCh
		binary.Write(file, binary.NativeEndian, buf)
		binary.Write(file, binary.NativeEndian, newLineB)
		//fmt.Printf("%b\n", buf)
	}
}

func Load() {
	filePath := getCurrentLogFile()
	exist, err := util.PathExists(filePath)
	if err != nil {
		l.Panic().Msg("[X]load binlog failed")
	}
	if !exist {
		return
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		l.Panic().Msg("[X]load binlog failed")
	}
	scanner := bufio.NewScanner(file)
	// resize scanner's capacity for lines over 64K, see next example
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		msg := scanner.Bytes()
		last := msg[len(msg)-1]
		if last == newLineB { // remove last '\n'
			msg = msg[:len(msg)-1]
		}
		msg = msg[4:] // remove msg length header

		icmd, _, err := proto.DecodeHeader(msg)
		if err != nil {
			l.Err(err)
		}

		if len(msg) < 2 {
			l.Warn().Msgf("[X]invalid binlog: %b", msg)
			continue
			//return resp.ErrorResponse(errors.New(resp.RES_SYNTAX_ERROR))
		}

		var req datatype.Req
		err = json.Unmarshal(msg[2:], &req) // 2= 1 icmd and 1 ver
		if err != nil {
			l.Warn().Msgf("[X]invalid binlog: %b", msg)
			continue
			//return resp.ErrorResponse(err)
		}
		executor.Execute1(icmd, req)
	}

	if err := scanner.Err(); err != nil {
		l.Panic().Msg(err.Error())
	}
}

func getCurrentLogFile() string {
	return util.GetMainPath() + consts.SLASH + consts.DATA_DIR + consts.SLASH + "bin_000001.log"
}
