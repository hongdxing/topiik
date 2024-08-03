/*
* author: duan hongxing
* date: 2 Aug, 2024
* desc:
*	Sync log data from Partition Leader
 */

package persistence

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
	"topiik/cluster"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/resp"
)

var ticker *time.Ticker
var ptnLeaderAddr string

func Sync() {
	ticker = time.NewTicker(1 * time.Second)

	for {
		<-ticker.C
		if cluster.GetNodeInfo().Id == ptnLeaderAddr { // if current node is Leader, do nothing
			break
		}
		doSync()
	}

}

func doSync() {
	if len(ptnLeaderAddr) == 0 {
		getPartitionLeader()
	}
}

/*
* Get Partition Leader from Controller
 */
func getPartitionLeader() {
	clAddr := cluster.GetNodeStatus().LeaderControllerAddr
	if len(clAddr) == 0 {
		return
	}

	hostPort, err := util.SplitAddress(clAddr)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return
	}

	conn, err := util.PreapareSocketClient(hostPort[0] + ":" + hostPort[2])
	if err != nil {
		return
	}
	defer conn.Close()

	var req []byte
	var byteBuf = new(bytes.Buffer) // int to byte byte buf
	// 1 bytes of command
	binary.Write(byteBuf, binary.LittleEndian, cluster.RPC_GET_PL)
	req = append(req, byteBuf.Bytes()...)
	req = append(req, []byte(cluster.GetNodeInfo().Id)...) // send current worker node id to get leader id

	// enocde
	req, err = proto.EncodeB(req)
	if err != nil {
		l.Err(err).Msg(err.Error())
	}

	// send
	_, err = conn.Write(req)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return
	}

	reader := bufio.NewReader(conn)
	buf, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			l.Err(err).Msgf("sync::getPartitionLeader %s", err.Error())
		}
	}
	if len(buf) < resp.RESPONSE_HEADER_SIZE {
		l.Warn().Msgf("sync::getPartitionLeader invalid response len%v", len(buf))
	}

	// read response flag
	byteBuf.Reset()
	bufSlice := buf[4:5]
	byteBuf = bytes.NewBuffer(bufSlice)
	var flag int8
	err = binary.Read(byteBuf, binary.LittleEndian, &flag)
	if err != nil {
		fmt.Println("(err):")
	}
	if flag == 1 {
		ptnLeaderAddr = string(buf[resp.RESPONSE_HEADER_SIZE:])
	}
	l.Info().Msgf("sync::getPartitionLeader new ptnleaderAddr: %s", ptnLeaderAddr)
}
